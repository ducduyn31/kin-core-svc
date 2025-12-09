package conversation

import (
	"time"

	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type Participant struct {
	ID             uuid.UUID  `json:"id"`
	ConversationID uuid.UUID  `json:"conversation_id"`
	UserID         uuid.UUID  `json:"user_id"`
	IsMuted        bool       `json:"is_muted"`
	IsArchived     bool       `json:"is_archived"`
	LastReadAt     *time.Time `json:"last_read_at,omitempty"`
	JoinedAt       time.Time  `json:"joined_at"`
	LeftAt         *time.Time `json:"left_at,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func NewParticipant(conversationID, userID uuid.UUID) *Participant {
	now := time.Now()
	return &Participant{
		ID:             uid.New(),
		ConversationID: conversationID,
		UserID:         userID,
		IsMuted:        false,
		IsArchived:     false,
		JoinedAt:       now,
		UpdatedAt:      now,
	}
}

func (p *Participant) Mute() {
	p.IsMuted = true
	p.UpdatedAt = time.Now()
}

func (p *Participant) Unmute() {
	p.IsMuted = false
	p.UpdatedAt = time.Now()
}

func (p *Participant) Archive() {
	p.IsArchived = true
	p.UpdatedAt = time.Now()
}

func (p *Participant) Unarchive() {
	p.IsArchived = false
	p.UpdatedAt = time.Now()
}

func (p *Participant) MarkAsRead(at time.Time) {
	p.LastReadAt = &at
	p.UpdatedAt = time.Now()
}

func (p *Participant) Leave() {
	now := time.Now()
	p.LeftAt = &now
	p.UpdatedAt = now
}

func (p *Participant) IsActive() bool {
	return p.LeftAt == nil
}
