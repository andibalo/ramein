package pubsub

import (
	"github.com/andibalo/ramein/commons/rabbitmq"
	"go.uber.org/zap"
)

func (p *pubsub) LogPayload(topic string, m rabbitmq.Message) {
	p.Config.Logger().Info("Received topic : ", zap.String("topic", topic))
	p.Config.Logger().Info("payload : ", zap.Any("payload", m))
}
