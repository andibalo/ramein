package pubsub

import (
	"context"
	"github.com/andibalo/ramein/commons/rabbitmq"
	"github.com/andibalo/ramein/orion/internal/config"
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

	_ = p.Rmq.Subscribe(rabbitmq.SubscriberConfig{
		Topic:   CORE_NEW_USER_REGISTERED,
		Channel: p.Config.RabbitMQChannel(),
	}, p.CoreNewUserRegisteredHandler)
}

func (p *pubsub) CoreNewUserRegisteredHandler(c context.Context, message rabbitmq.Message) error {
	p.LogPayload(CORE_NEW_USER_REGISTERED, message)

	return nil
}
