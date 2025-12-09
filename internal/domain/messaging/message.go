package messaging

import (
	"time"

	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID  `json:"id"`
	ConversationID uuid.UUID  `json:"conversation_id"`
	SenderID       uuid.UUID  `json:"sender_id"`
	Content        Content    `json:"content"`
	ReplyToID      *uuid.UUID `json:"reply_to_id,omitempty"`
	IsEdited       bool       `json:"is_edited"`
	EditedAt       *time.Time `json:"edited_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	DeletedForAll  bool       `json:"-"`
}

func NewMessage(conversationID, senderID uuid.UUID, content Content) *Message {
	return &Message{
		ID:             uid.New(),
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		IsEdited:       false,
		CreatedAt:      time.Now(),
	}
}

func (m *Message) SetReplyTo(messageID uuid.UUID) {
	m.ReplyToID = &messageID
}

func (m *Message) Edit(content Content) {
	m.Content = content
	m.IsEdited = true
	now := time.Now()
	m.EditedAt = &now
}

func (m *Message) DeleteForMe() {
}

func (m *Message) DeleteForAll() {
	m.DeletedForAll = true
}

func (m *Message) IsDeleted() bool {
	return m.DeletedForAll
}

func (m *Message) CanEdit(editWindowMinutes int) bool {
	if m.Content.Type != ContentTypeText {
		return false
	}
	if m.IsDeleted() {
		return false
	}
	return time.Since(m.CreatedAt).Minutes() <= float64(editWindowMinutes)
}

type MessageDeletion struct {
	ID        uuid.UUID `json:"id"`
	MessageID uuid.UUID `json:"message_id"`
	UserID    uuid.UUID `json:"user_id"`
	DeletedAt time.Time `json:"deleted_at"`
}

func NewMessageDeletion(messageID, userID uuid.UUID) *MessageDeletion {
	return &MessageDeletion{
		ID:        uid.New(),
		MessageID: messageID,
		UserID:    userID,
		DeletedAt: time.Now(),
	}
}
