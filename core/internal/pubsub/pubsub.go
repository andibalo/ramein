package pubsub

import (
	pubsubCommons "github.com/andibalo/ramein/commons/pubsub"
	"github.com/andibalo/ramein/commons/rabbitmq"
	"github.com/andibalo/ramein/core/internal/config"
)

type PubSub interface {
	PublishNewUserRegistered(payload pubsubCommons.CoreNewRegisteredUserPayload) error
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

}
