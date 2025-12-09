package media

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, media *Media) error
	GetByID(ctx context.Context, id uuid.UUID) (*Media, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*Media, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Media, error)
	CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)
}

type Storage interface {
	Upload(ctx context.Context, key string, reader io.Reader, contentType string, size int64) (string, error)
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	GetURL(ctx context.Context, key string) (string, error)
	GetSignedURL(ctx context.Context, key string, expiresIn int) (string, error)
}
