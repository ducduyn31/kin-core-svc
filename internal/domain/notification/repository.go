package notification

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, notification *Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*Notification, error)
	Update(ctx context.Context, notification *Notification) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool, limit, offset int) ([]*Notification, error)
	CountUnread(ctx context.Context, userID uuid.UUID) (int64, error)
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	DeleteOlderThan(ctx context.Context, userID uuid.UUID, daysOld int) error

	CreatePreferences(ctx context.Context, prefs *NotificationPreferences) error
	GetPreferences(ctx context.Context, userID uuid.UUID) (*NotificationPreferences, error)
	UpdatePreferences(ctx context.Context, prefs *NotificationPreferences) error
}

type PushService interface {
	SendPush(ctx context.Context, token string, notification *Notification) error
	SendMultiplePush(ctx context.Context, tokens []string, notification *Notification) error
}
