package repository

import "github.com/andibalo/ramein/astra/internal/model"

type MessageRepository interface {
	GetMessagesByConversationID(conversationID string, page []byte, limit int) ([]*model.Message, error)
	SaveMessage(message model.Message) error
}
