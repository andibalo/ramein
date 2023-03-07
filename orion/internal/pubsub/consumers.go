package pubsub

import (
	"context"
	"github.com/andibalo/ramein/commons/rabbitmq"
)

func (p *pubsub) CoreNewUserRegisteredHandler(c context.Context, message rabbitmq.Message) error {
	p.LogPayload(CORE_NEW_USER_REGISTERED, message)

	return nil
}
