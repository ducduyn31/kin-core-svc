package presence

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Set(ctx context.Context, presence *Presence, ttl time.Duration) error
	Get(ctx context.Context, userID uuid.UUID) (*Presence, error)
	GetMultiple(ctx context.Context, userIDs []uuid.UUID) ([]*Presence, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	SetOnline(ctx context.Context, userID uuid.UUID, deviceType DeviceType, deviceID *string, ttl time.Duration) error
	SetOffline(ctx context.Context, userID uuid.UUID) error
	Heartbeat(ctx context.Context, userID uuid.UUID, ttl time.Duration) error

	SetActivity(ctx context.Context, activity *Activity) error
	GetActivity(ctx context.Context, userID uuid.UUID) (*Activity, error)
	ClearActivity(ctx context.Context, userID uuid.UUID) error

	SetTyping(ctx context.Context, indicator *TypingIndicator) error
	GetTypingUsers(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error)
	ClearTyping(ctx context.Context, userID, conversationID uuid.UUID) error

	SetPushToken(ctx context.Context, userID uuid.UUID, token string) error
	GetPushToken(ctx context.Context, userID uuid.UUID) (string, error)
	DeletePushToken(ctx context.Context, userID uuid.UUID) error
}
