package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

var (
	//ErrorConsumeTimeout if flag handler consume timeout and requeue
	ErrorConsumeTimeout = errors.New("Consume Timeout")
)

// SubscriberConfig subscriber high-level configuration
type SubscriberConfig struct {
	// Timeout in second
	Timeout int
	// Topic will subscribe
	Topic string
	// Channel is queue label
	// default using service name
	Channel string

	// HookOnUnmarshalFailed is hook function
	// When error unmarshal Message struct
	OnUnmarshalFailed func(OnUnmarshalFailed, error)
}

// PubSubService abstraction for publish and subscribe broker
type PubSubService interface {
	Publish(topic string, msg Message) error
	Subscribe(config SubscriberConfig, h HandlerSubscriber) error
	PublishWithContext(ctx context.Context, topic string, msg Message) error
	Close() error
}

// HandlerSubscriber function for subscriber
type HandlerSubscriber func(context.Context, Message) error

// Message is message format
type Message struct {
	// UUID is message id
	// Default uuid()
	UUID string `json:"uuid,omitempty"`
	// Timestamp unix on queue publish
	// Default now()
	Timestamp int `json:"timestamp,omitempty"`
	// Type can event type or etc
	Type string `json:"type,omitempty"`
	// Metadata is data to pass to header
	Metadata map[string]string `json:"metadata,omitempty"`
	// Payload is data
	Payload interface{} `json:"payload,omitempty"`
}

// OnUnmarshalFailed structure data on unmarshal errors
type OnUnmarshalFailed struct {
	// body is original data from broker
	Body []byte

	// UUID is request-id or message.uuid
	UUID string `json:"uuid,omitempty"`
}

// WorkerName is identifier is subscriber or publisher
type WorkerName string

const (
	WorkerPublisher  WorkerName = "publisher"
	WorkerSubscriber WorkerName = "subscriber"
)

// OnConnectionError is struct data on connection error
type OnConnectionError struct {
	Worker WorkerName
	Error  error
	Meta   map[string]interface{}
}

type rabbitmqConn struct {
	conn *amqp.Connection
	ch   map[string]*amqp.Channel
}
type rabitmqSvc struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	pubConn   *rabbitmqConn
	subConn   *rabbitmqConn
	config    RabitmqConfiguration
	sync.Mutex
}

// RabitmqConfiguration is configuration nsq broker
type RabitmqConfiguration struct {
	// AutoRecconnect is disable or enable
	// automatically reconnect
	// default true
	AutoRecconnect *bool

	// Max retry connection
	// With exponential backoff retry
	// when < 1 no max conn retry == infinite reconnect
	MaxConnRetries int

	// ExponentialMultiplier multiplier
	// recconnection delay, default 10s
	ExponentialMultiplier int

	// OnConnectionError is hook function
	// When connection error
	OnConnectionError func(OnConnectionError)

	// URL is rabitmq url
	URL string

	// If need enable disable service
	Enable bool

	//TODO: RMQ MIDDLEWARE
	// Factory middleware is func return new middleware
	//Middleware FactoryMiddleware
}

const (
	exchangeType = "topic"
	keyRequestId = "request-id"
)

func (ss *rabitmqSvc) consumeExecute(ctx context.Context, config SubscriberConfig, h HandlerSubscriber, d amqp.Delivery) {
	var err error
	// create value context timeout
	ctx, cancelFunc := context.WithTimeout(ctx, time.Duration(config.Timeout)*time.Second)

	// Update the context with the span for the subsequent reference.
	// ctx = opentracing.ContextWithSpan(ctx, sp)
	defer cancelFunc()

	// default
	exchangeName := "amq.topic"
	// parse domain topic
	splitTopic := strings.Split(config.Topic, ".")
	if len(splitTopic) > 1 {
		exchangeName = splitTopic[0]
	}
	queueName := fmt.Sprintf("%s.%s", config.Topic, config.Channel)

	//  if connection disconnected
	//  and auto reconnect enable, will resubscribe
	if ss.subConn.conn.IsClosed() {

		if ss.isAutoReconnect() {
			ss.Subscribe(config, h)
		} else {
			if ss.config.OnConnectionError != nil {
				ss.config.OnConnectionError(OnConnectionError{
					Error:  fmt.Errorf("connection rabbitmq disconnect"),
					Worker: WorkerSubscriber,
					Meta: map[string]interface{}{
						"topic":    config.Topic,
						"exchange": exchangeName,
						"queue":    queueName,
					},
				})
			}

		}
		return
	}

	msg := Message{}
	var requestId string
	if d.Headers != nil && d.Headers[keyRequestId] != nil {
		requestId = fmt.Sprintf("%s", d.Headers[keyRequestId])
	}

	if err := json.Unmarshal(d.Body, &msg); err != nil {
		if config.OnUnmarshalFailed != nil {
			config.OnUnmarshalFailed(OnUnmarshalFailed{
				Body: d.Body,
				UUID: requestId,
			}, err)
		}

		// no requeue
		d.Reject(false)
		log.Println(err)
		return
	}

	// if requestId header not
	// empty add it as msg.UUID
	if requestId != "" {
		msg.UUID = requestId
	}

	// inject Metadata from headers
	metaData := map[string]string{}
	for k, v := range d.Headers {
		if val, ok := v.(string); ok {
			metaData[k] = val
		}
	}
	msg.Metadata = metaData

	errChan := make(chan error)
	go func(ctx context.Context, msg Message) {
		errChan <- h(ctx, msg)
	}(ctx, msg)

	go func(ctx context.Context) {
		// wait timeout
		for {
			select {
			case <-ctx.Done():
				errChan <- ErrorConsumeTimeout
				return
			default:
				<-time.After(time.Duration(config.Timeout) * time.Second)
			}
		}
	}(ctx)

	// if error not nil wil requeue message
	err = <-errChan
	if err != nil {
		log.Println(err, msg.UUID)
		d.Reject(true)
		return
	}

	d.Ack(false)
}

// Subscribe is subscribe to rabitmq
func (ss *rabitmqSvc) Subscribe(config SubscriberConfig, h HandlerSubscriber) error {
	// default
	exchangeName := "amq.topic"
	// parse domain topic
	splitTopic := strings.Split(config.Topic, ".")
	if len(splitTopic) > 1 {
		exchangeName = splitTopic[0]
	}

	if config.Timeout < 1 {
		config.Timeout = 15
	}

	if config.Channel == "" {
		config.Channel = filepath.Base(os.Args[0])
	}

	attempt := 0
	var err error
	// prevent open multiple connections on same time
	ss.Lock()
stateConnect:
	// receive value error
	// before call stateConnect
	if err != nil {
		// no auto reconnect
		if !ss.isAutoReconnect() {
			return err
		}

		if ss.config.OnConnectionError != nil {
			ss.config.OnConnectionError(OnConnectionError{
				Error:  err,
				Worker: WorkerSubscriber,
				Meta: map[string]interface{}{
					"topic":    config.Topic,
					"exchange": exchangeName,
				},
			})
		}

		attempt++
		if ss.config.MaxConnRetries > 0 && attempt > ss.config.MaxConnRetries {
			log.Println("max retry exceeded", attempt)
			return err
		}

		// previous delay + current delay
		delay := time.Second * time.Duration(((attempt-1)+attempt)*ss.config.ExponentialMultiplier)
		log.Println(err, "reconnect with delay", delay, "second", ", attempt", attempt)
		<-time.After(delay)

		// clean last connection
		ss.subConn = nil
	}

	if ss.subConn == nil {
		var conn *amqp.Connection
		conn, err = amqp.Dial(ss.config.URL)
		if err != nil {
			goto stateConnect
		}

		ss.subConn = &rabbitmqConn{
			conn: conn,
			ch:   make(map[string]*amqp.Channel),
		}
	}

	var ch *amqp.Channel
	ch, err = ss.subConn.conn.Channel()
	if err != nil {
		goto stateConnect
	}

	// 1 channel per subscriber thread safe
	ss.subConn.ch[uuid.NewString()] = ch

	err = ch.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		goto stateConnect
	}
	queueName := fmt.Sprintf("%s.%s", config.Topic, config.Channel)
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		goto stateConnect
	}

	err = ch.QueueBind(
		q.Name,       // queue name
		config.Topic, // routing key
		exchangeName, // exchange
		false,
		nil,
	)

	if err != nil {
		goto stateConnect
	}

	msgs, err := ch.Consume(
		q.Name,         // queue
		config.Channel, // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	ss.Unlock()

	go func() {
		defer ch.Close()
		for {
			select {
			case <-ss.ctx.Done():
				return
			case d := <-msgs:
				//TODO: RMQ MIDDLEWARE
				// declare next middleware consumer function
				//var next MiddlewareNext = func(ctx context.Context) {
				//	ss.consumeExecute(ctx, config, h, d)
				//}

				ss.consumeExecute(context.Background(), config, h, d)

				//TODO: RMQ MIDDLEWARE
				// Generate middleware from factory function
				// and execute subscriber middleware
				//ss.config.Middleware(ss.ctx, config.Topic).MiddlewareSubscriber(ss.ctx, RabbitMQDelivery(d), next)
			}
		}
	}()

	return nil
}

// PublishWithContext is publish with passing context
// implement middleware
func (ss *rabitmqSvc) PublishWithContext(ctx context.Context, topic string, msg Message) error {
	var err error

	//TODO: RMQ MIDDLEWARE
	// declare next middleware function
	//var next MiddlewareNext = func(newCtx context.Context) {
	//	err = ss.publishWithContext(newCtx, topic, msg)
	//}
	// Generate middleware from factory function
	// and execute publisher middleware
	//ss.config.Middleware(ss.ctx, topic).MiddlewarePublisher(ctx, msg, next)

	err = ss.publishWithContext(ctx, topic, msg)
	return err
}

func (ss *rabitmqSvc) publishWithContext(ctx context.Context, topic string, msg Message) error {
	// default
	exchangeName := "amq.topic"
	// parse domain topic
	splitTopic := strings.Split(topic, ".")
	if len(splitTopic) > 1 {
		exchangeName = splitTopic[0]
	}

	// validation
	if msg.Payload == nil {
		return fmt.Errorf("Payload required")
	}

	// default now
	if msg.Timestamp < 1 {
		msg.Timestamp = int(toUnixMilliSecond(time.Now()))
	}

	// default uuid
	if msg.UUID == "" {
		msg.UUID = uuid.New().String()
	}

	// prefare headers for message
	// -- headers from payload
	headers := amqp.Table{keyRequestId: msg.UUID}
	for k, v := range msg.Metadata {
		headers[k] = v
	}

	//TODO: RMQ MIDDLEWARE
	// -- headers from ctx RabbitMQHeaderPrefix
	//headersTracing, ok := ctx.Value(RabbitMQHeaderPrefix).(MapOpenTracingCarrier)
	//if ok {
	//	for k, v := range headersTracing {
	//		headers[k] = v
	//	}
	//}

	// prefare body for message
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// init amqp publish data
	msgPublish := amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
		Headers:     headers,
	}

	attempt := 0
stateConnect:
	// receive value error
	// before call stateConnect
	if err != nil {
		// no auto reconnect
		if !ss.isAutoReconnect() {
			return err
		}

		// hook  when disconnected
		if ss.config.OnConnectionError != nil {
			ss.config.OnConnectionError(OnConnectionError{
				Error:  err,
				Worker: WorkerPublisher,
				Meta: map[string]interface{}{
					"topic":   topic,
					"uuid":    msg.UUID,
					"payload": msg.Payload,
				},
			})

		}

		attempt++
		if ss.config.MaxConnRetries > 0 && attempt > ss.config.MaxConnRetries {
			log.Println("max retry exceeded", attempt)
			return err
		}
		// previous delay + current delay
		delay := time.Second * time.Duration(((attempt-1)+attempt)*ss.config.ExponentialMultiplier)
		log.Println("reconnect with delay", delay, "second", ", attempt", attempt)
		<-time.After(delay)

		// clean last connection
		ss.pubConn = nil
	}

	if ss.pubConn == nil {
		var conn *amqp.Connection
		conn, err = amqp.Dial(ss.config.URL)
		if err != nil {
			goto stateConnect
		}

		ss.pubConn = &rabbitmqConn{
			conn: conn,
			ch:   make(map[string]*amqp.Channel),
		}
	}

	ch := ss.pubConn.ch["default"]
	if ch == nil {
		ss.pubConn.ch["default"], err = ss.pubConn.conn.Channel()
		if err != nil {
			return err
		}
		ch = ss.pubConn.ch["default"]
	}

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		// do reconnect when connection is closed
		if ss.pubConn.conn.IsClosed() {
			goto stateConnect
		}
		return err
	}

	err = ch.Publish(
		exchangeName, // exchange
		topic,        // routing key
		false,        // mandatory
		false,        // immediate
		msgPublish,
	)

	if err != nil {
		// do reconnect when connection is closed
		if ss.pubConn.conn.IsClosed() {
			goto stateConnect
		}
	}

	return err
}

// Publish is publish to rabitmq broker using app context
// payload & type is required
func (ss *rabitmqSvc) Publish(topic string, msg Message) error {
	return ss.PublishWithContext(ss.ctx, topic, msg)
}

// Close is close channels and connections
func (ss *rabitmqSvc) Close() error {
	ss.ctxCancel()

	if ss.pubConn != nil {
		for _, ch := range ss.pubConn.ch {
			ch.Close()
		}
		ss.pubConn.conn.Close()
	}

	if ss.subConn != nil {
		for _, ch := range ss.subConn.ch {
			ch.Close()
		}
		ss.subConn.conn.Close()
	}

	return nil
}

func (ss *rabitmqSvc) isAutoReconnect() bool {
	if ss.config.AutoRecconnect == nil {
		return true
	}

	return *ss.config.AutoRecconnect
}

// NewRabitmq new abstraction rabbitmq pubsub
func NewRabitmq(config RabitmqConfiguration) PubSubService {
	return NewRabitmqWithContext(context.Background(), config)
}

// NewRabitmqWithContext new abstraction rabbitmq pubsub
// with passing context
func NewRabitmqWithContext(ctx context.Context, config RabitmqConfiguration) PubSubService {

	if config.ExponentialMultiplier < 1 {
		config.ExponentialMultiplier = 10
	}

	//TODO: RMQ MIDDLEWARE
	// set default middleware
	//if config.Middleware == nil {
	//	config.Middleware = newRabbitMqMw
	//}

	ctx, ctxCancel := context.WithCancel(ctx)
	ss := &rabitmqSvc{
		config:    config,
		ctx:       ctx,
		ctxCancel: ctxCancel,
	}
	return ss
}
