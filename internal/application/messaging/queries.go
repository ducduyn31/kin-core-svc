package messaging

import (
	"time"

	"github.com/google/uuid"
)

type GetMessageQuery struct {
	MessageID uuid.UUID
	UserID    uuid.UUID // For permission check
}

type ListMessagesQuery struct {
	ConversationID uuid.UUID
	UserID         uuid.UUID // For permission check
	Cursor         *time.Time
	Limit          int
}

type ListMessageReactionsQuery struct {
	MessageID uuid.UUID
	UserID    uuid.UUID // For permission check
}

type SearchMessagesQuery struct {
	ConversationID uuid.UUID
	UserID         uuid.UUID // For permission check
	Query          string
	Limit          int
}

type GetUnreadCountQuery struct {
	ConversationID uuid.UUID
	UserID         uuid.UUID
	Since          time.Time
}
