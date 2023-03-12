package pubsub

import (
	pubsubCommons "github.com/andibalo/ramein/commons/pubsub"
	"github.com/andibalo/ramein/commons/rabbitmq"
	"go.uber.org/zap"
)

func (p *pubsub) PublishNewUserRegistered(payload pubsubCommons.CoreNewRegisteredUserPayload) error {

	msg := rabbitmq.Message{
		Type:    pubsubCommons.CORE_NEW_USER_REGISTERED,
		Payload: payload,
	}

	err := p.Rmq.Publish(pubsubCommons.CORE_NEW_USER_REGISTERED, msg)

	if err != nil {
		p.Config.Logger().Error("Error publishing to queue", zap.String("topic",
			pubsubCommons.CORE_NEW_USER_REGISTERED), zap.Any("payload", payload), zap.Error(err))

		return err
	}

	return nil
}
