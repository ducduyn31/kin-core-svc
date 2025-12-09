package location

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateOrUpdate(ctx context.Context, location *Location) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Location, error)
	GetByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]*Location, error)
	Delete(ctx context.Context, userID uuid.UUID) error

	CreateHistory(ctx context.Context, history *LocationHistory) error
	ListHistoryByUser(ctx context.Context, userID uuid.UUID, from, to time.Time, limit int) ([]*LocationHistory, error)
	DeleteHistoryOlderThan(ctx context.Context, userID uuid.UUID, before time.Time) error

	CreatePlace(ctx context.Context, place *Place) error
	GetPlaceByID(ctx context.Context, id uuid.UUID) (*Place, error)
	UpdatePlace(ctx context.Context, place *Place) error
	DeletePlace(ctx context.Context, id uuid.UUID) error
	ListPlacesByUser(ctx context.Context, userID uuid.UUID) ([]*Place, error)
	FindPlaceByLocation(ctx context.Context, userID uuid.UUID, lat, lng float64) (*Place, error)

	CreateCheckIn(ctx context.Context, checkIn *CheckIn) error
	GetCheckInByID(ctx context.Context, id uuid.UUID) (*CheckIn, error)
	ListCheckInsByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*CheckIn, error)
	ListCheckInsByPlace(ctx context.Context, placeID uuid.UUID, limit int) ([]*CheckIn, error)
	GetLatestCheckIn(ctx context.Context, userID uuid.UUID) (*CheckIn, error)
}
