package presence

import (
	"time"

	"github.com/google/uuid"
)

type ActivityType string

const (
	ActivityTypeIdle      ActivityType = "idle"
	ActivityTypeTyping    ActivityType = "typing"
	ActivityTypeViewing   ActivityType = "viewing"
	ActivityTypeInCall    ActivityType = "in_call"
	ActivityTypeRecording ActivityType = "recording"
)

type Activity struct {
	UserID         uuid.UUID    `json:"user_id"`
	Type           ActivityType `json:"type"`
	ConversationID *uuid.UUID   `json:"conversation_id,omitempty"` // Where the activity is happening
	Description    *string      `json:"description,omitempty"`
	StartedAt      time.Time    `json:"started_at"`
	ExpiresAt      time.Time    `json:"expires_at"` // When the activity auto-clears
}

func NewActivity(userID uuid.UUID, activityType ActivityType, duration time.Duration) *Activity {
	now := time.Now()
	return &Activity{
		UserID:    userID,
		Type:      activityType,
		StartedAt: now,
		ExpiresAt: now.Add(duration),
	}
}

func (a *Activity) SetConversation(conversationID uuid.UUID) {
	a.ConversationID = &conversationID
}

func (a *Activity) SetDescription(description *string) {
	a.Description = description
}

func (a *Activity) IsExpired() bool {
	return time.Now().After(a.ExpiresAt)
}

func (a *Activity) Refresh(duration time.Duration) {
	a.ExpiresAt = time.Now().Add(duration)
}

type TypingIndicator struct {
	UserID         uuid.UUID `json:"user_id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	StartedAt      time.Time `json:"started_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}

func NewTypingIndicator(userID, conversationID uuid.UUID, duration time.Duration) *TypingIndicator {
	now := time.Now()
	return &TypingIndicator{
		UserID:         userID,
		ConversationID: conversationID,
		StartedAt:      now,
		ExpiresAt:      now.Add(duration),
	}
}

func (t *TypingIndicator) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

func (t *TypingIndicator) Refresh(duration time.Duration) {
	t.ExpiresAt = time.Now().Add(duration)
}
