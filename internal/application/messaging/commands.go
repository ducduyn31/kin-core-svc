package messaging

import (
	"github.com/danielng/kin-core-svc/internal/domain/messaging"
	"github.com/google/uuid"
)

type SendMessageCommand struct {
	ConversationID uuid.UUID
	SenderID       uuid.UUID
	Content        messaging.Content
	ReplyToID      *uuid.UUID
}

type EditMessageCommand struct {
	MessageID uuid.UUID
	UserID    uuid.UUID
	Content   messaging.Content
}

type DeleteMessageCommand struct {
	MessageID   uuid.UUID
	UserID      uuid.UUID
	ForEveryone bool
}

type AddReactionCommand struct {
	MessageID uuid.UUID
	UserID    uuid.UUID
	Emoji     string
}

type RemoveReactionCommand struct {
	MessageID uuid.UUID
	UserID    uuid.UUID
	Emoji     string
}

type MarkAsReadCommand struct {
	ConversationID uuid.UUID
	UserID         uuid.UUID
	UpToMessageID  uuid.UUID
}

type MarkAsDeliveredCommand struct {
	MessageIDs []uuid.UUID
	UserID     uuid.UUID
}
