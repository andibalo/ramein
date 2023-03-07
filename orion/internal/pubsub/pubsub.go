package pubsub

import (
	"github.com/andibalo/ramein/commons/rabbitmq"
	"github.com/andibalo/ramein/orion/internal/config"
	"go.uber.org/zap"
)

type PubSub interface {
}

type pubsub struct {
	Config config.Config
	Rmq    rabbitmq.PubSubService
}

func NewPubSub(cfg config.Config, rmq rabbitmq.PubSubService) *pubsub {

	return &pubsub{
		Config: cfg,
		Rmq:    rmq,
	}
}

func (p *pubsub) InitSubscribers() {

	err := p.Rmq.Subscribe(rabbitmq.SubscriberConfig{
		Topic:   CORE_NEW_USER_REGISTERED,
		Channel: p.Config.RabbitMQChannel(),
	}, p.CoreNewUserRegisteredHandler)

	if err != nil {
		p.Config.Logger().Error("Error subscribing to topic", zap.String("topic", CORE_NEW_USER_REGISTERED))
	}

	p.Config.Logger().Info("Success subscribing to topic", zap.String("topic", CORE_NEW_USER_REGISTERED))
}
