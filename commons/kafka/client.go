package kafka

import "go.uber.org/zap"

type ClientOption func(*kafkaClient)

type Client interface {
	AddOptions(options ...ClientOption)
	GetSyncProducer() SyncProducer
}

func NewKafkaClient(options ...ClientOption) Client {
	client := &kafkaClient{}

	if len(options) > 0 {
		client.AddOptions(options...)
	}

	return client
}

func WithSyncProducer(v SyncProducer) ClientOption {
	return func(client *kafkaClient) {
		client.syncProducer = v
	}
}

type kafkaClient struct {
	logger       *zap.Logger
	syncProducer SyncProducer
}

func (client *kafkaClient) AddOptions(optionFuncs ...ClientOption) {
	for _, optionFunc := range optionFuncs {
		optionFunc(client)
	}
}

func (client *kafkaClient) GetSyncProducer() SyncProducer {
	return client.syncProducer
}
