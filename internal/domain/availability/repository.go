package availability

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateOrUpdate(ctx context.Context, availability *Availability) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Availability, error)
	GetByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]*Availability, error)

	CreateWindow(ctx context.Context, window *Window) error
	GetWindowByID(ctx context.Context, id uuid.UUID) (*Window, error)
	UpdateWindow(ctx context.Context, window *Window) error
	DeleteWindow(ctx context.Context, id uuid.UUID) error
	ListWindowsByUser(ctx context.Context, userID uuid.UUID) ([]*Window, error)
	ListActiveWindowsByUser(ctx context.Context, userID uuid.UUID) ([]*Window, error)
	ListWindowsByUserAndWeekday(ctx context.Context, userID uuid.UUID, weekday Weekday) ([]*Window, error)

	CreateAutoRule(ctx context.Context, rule *AutoRule) error
	GetAutoRuleByID(ctx context.Context, id uuid.UUID) (*AutoRule, error)
	UpdateAutoRule(ctx context.Context, rule *AutoRule) error
	DeleteAutoRule(ctx context.Context, id uuid.UUID) error
	ListAutoRulesByUser(ctx context.Context, userID uuid.UUID) ([]*AutoRule, error)
	ListActiveAutoRulesByUser(ctx context.Context, userID uuid.UUID) ([]*AutoRule, error)
}
