package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByAuth0Sub(ctx context.Context, auth0Sub string) (*User, error)
	Update(ctx context.Context, user *User) error
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*User, error)
	SearchByDisplayName(ctx context.Context, query string, limit int) ([]*User, error)

	CreatePreferences(ctx context.Context, prefs *Preferences) error
	GetPreferences(ctx context.Context, userID uuid.UUID) (*Preferences, error)
	UpdatePreferences(ctx context.Context, prefs *Preferences) error
}
