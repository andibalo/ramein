package kafka

import (
	"context"
	kafkago "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type SyncProducerOption func(*syncProducer)
type CloseProducerFunc func() error

type SyncProducer interface {
	Publish(key string, value []byte)
	Close() error
}

func NewSyncProducer(brokerHosts []string, topic string, options ...SyncProducerOption) (SyncProducer, error) {

	w := kafkago.Writer{
		Addr:  kafkago.TCP(brokerHosts...),
		Topic: topic,
	}

	sp := &syncProducer{
		producer: &w,
	}

	if len(options) > 0 {
		sp.AddOptions(options...)
	}

	return sp, nil
}

func WithLogger(v *zap.Logger) SyncProducerOption {
	return func(client *syncProducer) {
		client.logger = v
	}
}

type syncProducer struct {
	logger   *zap.Logger
	producer *kafkago.Writer
}

func (sp *syncProducer) AddOptions(optionFuncs ...SyncProducerOption) {
	for _, optionFunc := range optionFuncs {
		optionFunc(sp)
	}
}

func (sp *syncProducer) Close() error {
	return sp.producer.Close()
}

func (sp *syncProducer) Publish(key string, value []byte) {
	message := kafkago.Message{
		Key:   []byte(key),
		Value: value,
	}

	err := sp.producer.WriteMessages(context.Background(), message)
	if err != nil {
		sp.logger.Error("unable to produce kafka message", zap.Error(err))
	} else {
		sp.logger.Info("succesfully publish message to topic", zap.String("topic", sp.producer.Topic))
	}
}
