package service

type MessageService interface {
	GetMessagesByConversationID(conversationID string)
}
