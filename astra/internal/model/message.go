package model

import (
	"github.com/gocql/gocql"
	"time"
)

type Message struct {
	ConversationID    gocql.UUID
	MessageID         gocql.UUID
	ConversationName  string
	FromUserID        string
	FromUserNumber    string
	FromUserFirstName string
	FromUserLastName  string
	FromUserEmail     string
	TextContent       string
	SentAt            time.Time
	SeenAt            time.Time
	CreatedBy         string
	CreatedAt         time.Time
	UpdatedBy         string
	UpdatedAt         *time.Time
	DeletedBy         string
	DeletedAt         *time.Time
}
