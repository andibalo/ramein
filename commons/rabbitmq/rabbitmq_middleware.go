package pubsub

//TODO : rabbit mq middleware
// FactoryMiddleware is factory function to return
// PubSubMiddleware middleware
//type FactoryMiddleware func(ctx context.Context, topic string) PubSubMiddleware
//
//// MiddlewareNext is function next middleware / action
//type MiddlewareNext func(context.Context)
//
//// PubSubMiddleware interface middleware
//type PubSubMiddleware interface {
//	MiddlewarePublisher(ctx context.Context, msg Message, next MiddlewareNext)
//	MiddlewareSubscriber(ctx context.Context, delivery interface{}, next MiddlewareNext)
//}

//// RabbitMQHeaderPrefix is context key for RabbitMQ header
//const RabbitMQHeaderPrefix = "HEADER:RabbitMQ"
//
//type RabbitMQDelivery amqp.Delivery
//
//type rabbitMqMw struct {
//	ctx   context.Context
//	topic string
//}
//
//// MiddlewarePublisher is middleware on publish message
//// next is continue middleware / do publish
//// message is message will be publish
//func (hook *rabbitMqMw) MiddlewarePublisher(ctx context.Context, msg Message, next MiddlewareNext) {
//
//	var (
//		span opentracing.Span
//	)
//
//	// start new span or continue span from context
//	span, ctx = opentracing.StartSpanFromContext(ctx, "publisher."+hook.topic)
//	defer span.Finish()
//
//	// Inject the span context into the AMQP header.
//	data := make(MapOpenTracingCarrier)
//	// source https://github.com/opentracing-contrib/go-amqp/blob/master/amqptracer/tracer.go#L30
//	err := span.Tracer().Inject(span.Context(), opentracing.TextMap, data)
//	if err != nil {
//		log.Println("rabbitmq error inject span context", err)
//	}
//
//	// inject rabbitMQ header from context
//	ctx = context.WithValue(ctx, RabbitMQHeaderPrefix, data)
//
//	next(ctx)
//}
//
//// MiddlewareSubscriber is middleware on consume message
//// next is continue middleware / do consume
//// delivery can be any but currently only support amqp delivery
//func (hook *rabbitMqMw) MiddlewareSubscriber(ctx context.Context, delivery interface{}, next MiddlewareNext) {
//	d, ok := delivery.(RabbitMQDelivery)
//	if !ok {
//		next(ctx)
//		return
//	}
//	// extract span from rabbitMQ header using opentracing
//	// continue span from traceparent or init new tracing when have no traceparent header
//	// source https://github.com/opentracing-contrib/go-amqp/blob/master/amqptracer/tracer.go#L54
//	spCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, MapOpenTracingCarrier(d.Headers))
//	if err != nil {
//		log.Println("error extracting", err)
//	}
//
//	// Extract the span context out of the AMQP header.
//	sp := opentracing.StartSpan(
//		"consumer."+hook.topic,
//		opentracing.ChildOf(spCtx),
//	)
//	defer sp.Finish()
//
//	// Update the context with the span for the subsequent reference.
//	ctx = opentracing.ContextWithSpan(ctx, sp)
//
//	next(ctx)
//
//}
//
//func newRabbitMqMw(ctx context.Context, topic string) PubSubMiddleware {
//	if !opentracing.IsGlobalTracerRegistered() {
//		log.Println("opentracing global tracer  is not registered")
//	}
//	return &rabbitMqMw{ctx: ctx, topic: topic}
//}
