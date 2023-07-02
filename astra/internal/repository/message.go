package repository

import (
	"fmt"
	"github.com/andibalo/ramein/astra/internal/model"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/table"
	"go.uber.org/zap"
)

type MessageRepo struct {
	MessageTable *table.Table
	dbSession    gocqlx.Session
	logger       *zap.Logger
}

func NewMessageRepository(dbSession gocqlx.Session, logger *zap.Logger) MessageRepository {

	messageMetadata := table.Metadata{
		Name: "message_by_conversation_id",
		Columns: []string{
			"conversation_id",
			"message_id",
			"conversation_name",
			"from_user_id",
			"from_user_number",
			"from_user_first_name",
			"from_user_last_name",
			"from_user_email",
			"text_content",
			"sent_at",
			"created_by",
			"created_at",
			"updated_by",
			"updated_at",
			"deleted_by",
			"deleted_at",
		},
		PartKey: []string{"conversation_id"},
		SortKey: []string{"sent_at"},
	}

	messageTable := table.New(messageMetadata)

	return &MessageRepo{
		MessageTable: messageTable,
		dbSession:    dbSession,
		logger:       logger,
	}
}

func (r *MessageRepo) GetMessagesByConversationID(conversationID string, page []byte, limit int) ([]*model.Message, error) {
	var (
		messages []*model.Message
		err      error
	)

	q := r.MessageTable.SelectQuery(r.dbSession).Bind(conversationID)
	defer q.Release()
	q.PageState(page)
	q.PageSize(limit)

	iter := q.Iter()

	err = iter.Select(&messages)
	if err != nil {
		panic(err)
	}

	fmt.Println(iter.PageState())

	return messages, nil
}

func (r *MessageRepo) SaveMessage(message model.Message) error {

	insertMessage := r.MessageTable.InsertQuery(r.dbSession)

	insertMessage.BindStruct(message)
	if err := insertMessage.ExecRelease(); err != nil {
		r.logger.Error("error inserting message to db", zap.Error(err))

		return err
	}

	return nil
}
