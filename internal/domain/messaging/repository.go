package messaging

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, message *Message) error
	GetByID(ctx context.Context, id uuid.UUID) (*Message, error)
	Update(ctx context.Context, message *Message) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByConversation(ctx context.Context, conversationID uuid.UUID, cursor *time.Time, limit int) ([]*Message, error)
	ListByConversationAfter(ctx context.Context, conversationID uuid.UUID, afterID uuid.UUID, limit int) ([]*Message, error)
	CountByConversation(ctx context.Context, conversationID uuid.UUID) (int64, error)
	CountUnreadByUser(ctx context.Context, conversationID, userID uuid.UUID, after time.Time) (int64, error)
	GetLatestByConversation(ctx context.Context, conversationID uuid.UUID) (*Message, error)
	SearchInConversation(ctx context.Context, conversationID uuid.UUID, query string, limit int) ([]*Message, error)

	CreateDeletion(ctx context.Context, deletion *MessageDeletion) error
	GetDeletion(ctx context.Context, messageID, userID uuid.UUID) (*MessageDeletion, error)
	ListDeletedByUser(ctx context.Context, userID uuid.UUID, conversationID uuid.UUID) ([]uuid.UUID, error)

	CreateReceipt(ctx context.Context, receipt *Receipt) error
	GetReceipt(ctx context.Context, messageID, userID uuid.UUID) (*Receipt, error)
	UpdateReceipt(ctx context.Context, receipt *Receipt) error
	ListReceiptsByMessage(ctx context.Context, messageID uuid.UUID) ([]*Receipt, error)
	BulkUpdateReceiptsDelivered(ctx context.Context, messageIDs []uuid.UUID, userID uuid.UUID) error
	BulkUpdateReceiptsRead(ctx context.Context, messageIDs []uuid.UUID, userID uuid.UUID) error

	CreateReaction(ctx context.Context, reaction *Reaction) error
	GetReaction(ctx context.Context, messageID, userID uuid.UUID, emoji string) (*Reaction, error)
	DeleteReaction(ctx context.Context, id uuid.UUID) error
	DeleteUserReaction(ctx context.Context, messageID, userID uuid.UUID, emoji string) error
	ListReactionsByMessage(ctx context.Context, messageID uuid.UUID) ([]*Reaction, error)
	CountReactionsByMessage(ctx context.Context, messageID uuid.UUID) (map[string]int, error)
}
