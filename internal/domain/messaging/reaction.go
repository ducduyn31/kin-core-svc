package messaging

import (
	"time"

	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type Reaction struct {
	ID        uuid.UUID `json:"id"`
	MessageID uuid.UUID `json:"message_id"`
	UserID    uuid.UUID `json:"user_id"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
}

func NewReaction(messageID, userID uuid.UUID, emoji string) *Reaction {
	return &Reaction{
		ID:        uid.New(),
		MessageID: messageID,
		UserID:    userID,
		Emoji:     emoji,
		CreatedAt: time.Now(),
	}
}

type ReactionSummary struct {
	Emoji string      `json:"emoji"`
	Count int         `json:"count"`
	Users []uuid.UUID `json:"users,omitempty"`
}

func GetReactionSummaries(reactions []*Reaction) []ReactionSummary {
	emojiMap := make(map[string][]uuid.UUID)
	order := make([]string, 0)

	for _, r := range reactions {
		if _, exists := emojiMap[r.Emoji]; !exists {
			order = append(order, r.Emoji)
		}
		emojiMap[r.Emoji] = append(emojiMap[r.Emoji], r.UserID)
	}

	summaries := make([]ReactionSummary, 0, len(emojiMap))
	for _, emoji := range order {
		users := emojiMap[emoji]
		summaries = append(summaries, ReactionSummary{
			Emoji: emoji,
			Count: len(users),
			Users: users,
		})
	}

	return summaries
}
