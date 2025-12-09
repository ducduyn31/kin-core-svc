package conversation

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, conversation *Conversation) error
	GetByID(ctx context.Context, id uuid.UUID) (*Conversation, error)
	GetDirectConversation(ctx context.Context, userID1, userID2 uuid.UUID) (*Conversation, error)
	GetByCircleID(ctx context.Context, circleID uuid.UUID) (*Conversation, error)
	Update(ctx context.Context, conversation *Conversation) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByUser(ctx context.Context, userID uuid.UUID, includeArchived bool, limit, offset int) ([]*Conversation, error)
	CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)

	AddParticipant(ctx context.Context, participant *Participant) error
	GetParticipant(ctx context.Context, conversationID, userID uuid.UUID) (*Participant, error)
	UpdateParticipant(ctx context.Context, participant *Participant) error
	RemoveParticipant(ctx context.Context, conversationID, userID uuid.UUID) error
	ListParticipants(ctx context.Context, conversationID uuid.UUID) ([]*Participant, error)
	ListActiveParticipants(ctx context.Context, conversationID uuid.UUID) ([]*Participant, error)
	IsParticipant(ctx context.Context, conversationID, userID uuid.UUID) (bool, error)
	CountUnreadConversations(ctx context.Context, userID uuid.UUID) (int64, error)
}
