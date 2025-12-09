package contact

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, contact *Contact) error
	GetByID(ctx context.Context, id uuid.UUID) (*Contact, error)
	GetByUserAndContact(ctx context.Context, userID, contactID uuid.UUID) (*Contact, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Contact, error)
	ListFavorites(ctx context.Context, userID uuid.UUID) ([]*Contact, error)
	ListBlocked(ctx context.Context, userID uuid.UUID) ([]*Contact, error)
	Update(ctx context.Context, contact *Contact) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)

	CreateRequest(ctx context.Context, request *ContactRequest) error
	GetRequestByID(ctx context.Context, id uuid.UUID) (*ContactRequest, error)
	GetPendingRequest(ctx context.Context, fromUserID, toUserID uuid.UUID) (*ContactRequest, error)
	ListPendingRequestsForUser(ctx context.Context, userID uuid.UUID) ([]*ContactRequest, error)
	ListSentRequests(ctx context.Context, userID uuid.UUID) ([]*ContactRequest, error)
	UpdateRequest(ctx context.Context, request *ContactRequest) error
}
