package conversation

import (
	"time"

	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type ConversationType string

const (
	ConversationTypeDirect ConversationType = "direct" // One-on-one chat
	ConversationTypeCircle ConversationType = "circle" // Circle group chat
)

type Conversation struct {
	ID            uuid.UUID        `json:"id"`
	Type          ConversationType `json:"type"`
	CircleID      *uuid.UUID       `json:"circle_id,omitempty"` // Only for circle conversations
	Name          *string          `json:"name,omitempty"`      // Optional name for group chats
	Avatar        *string          `json:"avatar,omitempty"`
	LastMessageID *uuid.UUID       `json:"last_message_id,omitempty"`
	LastMessageAt *time.Time       `json:"last_message_at,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

func NewDirectConversation() *Conversation {
	now := time.Now()
	return &Conversation{
		ID:        uid.New(),
		Type:      ConversationTypeDirect,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewCircleConversation(circleID uuid.UUID, name *string) *Conversation {
	now := time.Now()
	return &Conversation{
		ID:        uid.New(),
		Type:      ConversationTypeCircle,
		CircleID:  &circleID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (c *Conversation) UpdateLastMessage(messageID uuid.UUID) {
	now := time.Now()
	c.LastMessageID = &messageID
	c.LastMessageAt = &now
	c.UpdatedAt = now
}

func (c *Conversation) SetName(name *string) {
	c.Name = name
	c.UpdatedAt = time.Now()
}

func (c *Conversation) SetAvatar(avatar *string) {
	c.Avatar = avatar
	c.UpdatedAt = time.Now()
}

func (c *Conversation) IsDirect() bool {
	return c.Type == ConversationTypeDirect
}

func (c *Conversation) IsCircle() bool {
	return c.Type == ConversationTypeCircle
}
