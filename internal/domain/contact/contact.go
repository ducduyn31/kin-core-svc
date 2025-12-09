package contact

import (
	"time"

	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type Contact struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	ContactID  uuid.UUID `json:"contact_id"`
	Nickname   *string   `json:"nickname,omitempty"`
	IsFavorite bool      `json:"is_favorite"`
	IsBlocked  bool      `json:"is_blocked"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewContact(userID, contactID uuid.UUID) *Contact {
	now := time.Now()
	return &Contact{
		ID:         uid.New(),
		UserID:     userID,
		ContactID:  contactID,
		IsFavorite: false,
		IsBlocked:  false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (c *Contact) SetNickname(nickname *string) {
	c.Nickname = nickname
	c.UpdatedAt = time.Now()
}

func (c *Contact) SetFavorite(favorite bool) {
	c.IsFavorite = favorite
	c.UpdatedAt = time.Now()
}

func (c *Contact) Block() {
	c.IsBlocked = true
	c.UpdatedAt = time.Now()
}

func (c *Contact) Unblock() {
	c.IsBlocked = false
	c.UpdatedAt = time.Now()
}

type ContactRequest struct {
	ID         uuid.UUID            `json:"id"`
	FromUserID uuid.UUID            `json:"from_user_id"`
	ToUserID   uuid.UUID            `json:"to_user_id"`
	Message    *string              `json:"message,omitempty"`
	Status     ContactRequestStatus `json:"status"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
}

type ContactRequestStatus string

const (
	ContactRequestStatusPending  ContactRequestStatus = "pending"
	ContactRequestStatusAccepted ContactRequestStatus = "accepted"
	ContactRequestStatusRejected ContactRequestStatus = "rejected"
)

func NewContactRequest(fromUserID, toUserID uuid.UUID, message *string) *ContactRequest {
	now := time.Now()
	return &ContactRequest{
		ID:         uid.New(),
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Message:    message,
		Status:     ContactRequestStatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (cr *ContactRequest) Accept() {
	cr.Status = ContactRequestStatusAccepted
	cr.UpdatedAt = time.Now()
}

func (cr *ContactRequest) Reject() {
	cr.Status = ContactRequestStatusRejected
	cr.UpdatedAt = time.Now()
}

func (cr *ContactRequest) IsPending() bool {
	return cr.Status == ContactRequestStatusPending
}
