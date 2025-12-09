package messaging

import (
	"context"
	"log/slog"

	"github.com/danielng/kin-core-svc/internal/domain/conversation"
	"github.com/danielng/kin-core-svc/internal/domain/messaging"
	"github.com/google/uuid"
)

type Service struct {
	messageRepo      messaging.Repository
	conversationRepo conversation.Repository
	logger           *slog.Logger
	editWindowMins   int
}

func NewService(
	messageRepo messaging.Repository,
	conversationRepo conversation.Repository,
	logger *slog.Logger,
) *Service {
	return &Service{
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
		logger:           logger,
		editWindowMins:   15, // 15 minute edit window
	}
}

func (s *Service) SendMessage(ctx context.Context, cmd SendMessageCommand) (*messaging.Message, error) {
	isParticipant, err := s.conversationRepo.IsParticipant(ctx, cmd.ConversationID, cmd.SenderID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, conversation.ErrNotParticipant
	}

	if !messaging.IsValidContentType(cmd.Content.Type) {
		return nil, messaging.ErrInvalidContentType
	}

	msg := messaging.NewMessage(cmd.ConversationID, cmd.SenderID, cmd.Content)
	if cmd.ReplyToID != nil {
		msg.SetReplyTo(*cmd.ReplyToID)
	}

	if err := s.messageRepo.Create(ctx, msg); err != nil {
		s.logger.Error("failed to create message", "error", err)
		return nil, err
	}

	conv, err := s.conversationRepo.GetByID(ctx, cmd.ConversationID)
	if err == nil {
		conv.UpdateLastMessage(msg.ID)
		s.conversationRepo.Update(ctx, conv)
	}

	participants, err := s.conversationRepo.ListActiveParticipants(ctx, cmd.ConversationID)
	if err == nil {
		for _, p := range participants {
			if p.UserID != cmd.SenderID {
				receipt := messaging.NewReceipt(msg.ID, p.UserID)
				s.messageRepo.CreateReceipt(ctx, receipt)
			}
		}
	}

	s.logger.Info("message sent", "message_id", msg.ID, "conversation_id", cmd.ConversationID)
	return msg, nil
}

func (s *Service) GetMessage(ctx context.Context, query GetMessageQuery) (*messaging.Message, error) {
	msg, err := s.messageRepo.GetByID(ctx, query.MessageID)
	if err != nil {
		return nil, err
	}

	isParticipant, err := s.conversationRepo.IsParticipant(ctx, msg.ConversationID, query.UserID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, conversation.ErrNotParticipant
	}

	return msg, nil
}

func (s *Service) ListMessages(ctx context.Context, query ListMessagesQuery) ([]*messaging.Message, error) {
	isParticipant, err := s.conversationRepo.IsParticipant(ctx, query.ConversationID, query.UserID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, conversation.ErrNotParticipant
	}

	limit := query.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	return s.messageRepo.ListByConversation(ctx, query.ConversationID, query.Cursor, limit)
}

func (s *Service) EditMessage(ctx context.Context, cmd EditMessageCommand) (*messaging.Message, error) {
	msg, err := s.messageRepo.GetByID(ctx, cmd.MessageID)
	if err != nil {
		return nil, err
	}

	if msg.SenderID != cmd.UserID {
		return nil, messaging.ErrNotMessageSender
	}

	if !msg.CanEdit(s.editWindowMins) {
		return nil, messaging.ErrCannotEditMessage
	}

	msg.Edit(cmd.Content)

	if err := s.messageRepo.Update(ctx, msg); err != nil {
		s.logger.Error("failed to edit message", "error", err)
		return nil, err
	}

	return msg, nil
}

func (s *Service) DeleteMessage(ctx context.Context, cmd DeleteMessageCommand) error {
	msg, err := s.messageRepo.GetByID(ctx, cmd.MessageID)
	if err != nil {
		return err
	}

	if cmd.ForEveryone {
		if msg.SenderID != cmd.UserID {
			return messaging.ErrNotMessageSender
		}
		msg.DeleteForAll()
		if err := s.messageRepo.Update(ctx, msg); err != nil {
			s.logger.Error("failed to delete message for everyone", "error", err)
			return err
		}
	} else {
		deletion := messaging.NewMessageDeletion(cmd.MessageID, cmd.UserID)
		if err := s.messageRepo.CreateDeletion(ctx, deletion); err != nil {
			s.logger.Error("failed to delete message for user", "error", err)
			return err
		}
	}

	return nil
}

func (s *Service) AddReaction(ctx context.Context, cmd AddReactionCommand) (*messaging.Reaction, error) {
	msg, err := s.messageRepo.GetByID(ctx, cmd.MessageID)
	if err != nil {
		return nil, err
	}

	isParticipant, err := s.conversationRepo.IsParticipant(ctx, msg.ConversationID, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, conversation.ErrNotParticipant
	}

	existing, err := s.messageRepo.GetReaction(ctx, cmd.MessageID, cmd.UserID, cmd.Emoji)
	if err == nil && existing != nil {
		return existing, nil // Already reacted
	}

	reaction := messaging.NewReaction(cmd.MessageID, cmd.UserID, cmd.Emoji)
	if err := s.messageRepo.CreateReaction(ctx, reaction); err != nil {
		s.logger.Error("failed to add reaction", "error", err)
		return nil, err
	}

	return reaction, nil
}

func (s *Service) RemoveReaction(ctx context.Context, cmd RemoveReactionCommand) error {
	if err := s.messageRepo.DeleteUserReaction(ctx, cmd.MessageID, cmd.UserID, cmd.Emoji); err != nil {
		s.logger.Error("failed to remove reaction", "error", err)
		return err
	}
	return nil
}

func (s *Service) ListReactions(ctx context.Context, query ListMessageReactionsQuery) ([]*messaging.Reaction, error) {
	msg, err := s.messageRepo.GetByID(ctx, query.MessageID)
	if err != nil {
		return nil, err
	}

	isParticipant, err := s.conversationRepo.IsParticipant(ctx, msg.ConversationID, query.UserID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, conversation.ErrNotParticipant
	}

	return s.messageRepo.ListReactionsByMessage(ctx, query.MessageID)
}

func (s *Service) MarkAsRead(ctx context.Context, cmd MarkAsReadCommand) error {
	msg, err := s.messageRepo.GetByID(ctx, cmd.UpToMessageID)
	if err != nil {
		return err
	}

	messages, err := s.messageRepo.ListByConversation(ctx, cmd.ConversationID, &msg.CreatedAt, 1000)
	if err != nil {
		return err
	}

	messageIDs := make([]uuid.UUID, 0, len(messages))
	for _, m := range messages {
		if m.SenderID != cmd.UserID {
			messageIDs = append(messageIDs, m.ID)
		}
	}

	if len(messageIDs) > 0 {
		if err := s.messageRepo.BulkUpdateReceiptsRead(ctx, messageIDs, cmd.UserID); err != nil {
			s.logger.Error("failed to mark messages as read", "error", err)
			return err
		}
	}

	return nil
}

func (s *Service) MarkAsDelivered(ctx context.Context, cmd MarkAsDeliveredCommand) error {
	if len(cmd.MessageIDs) == 0 {
		return nil
	}

	if err := s.messageRepo.BulkUpdateReceiptsDelivered(ctx, cmd.MessageIDs, cmd.UserID); err != nil {
		s.logger.Error("failed to mark messages as delivered", "error", err)
		return err
	}

	return nil
}

func (s *Service) SearchMessages(ctx context.Context, query SearchMessagesQuery) ([]*messaging.Message, error) {
	isParticipant, err := s.conversationRepo.IsParticipant(ctx, query.ConversationID, query.UserID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, conversation.ErrNotParticipant
	}

	limit := query.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	return s.messageRepo.SearchInConversation(ctx, query.ConversationID, query.Query, limit)
}

func (s *Service) GetUnreadCount(ctx context.Context, query GetUnreadCountQuery) (int64, error) {
	return s.messageRepo.CountUnreadByUser(ctx, query.ConversationID, query.UserID, query.Since)
}
