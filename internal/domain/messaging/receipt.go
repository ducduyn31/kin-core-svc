package messaging

import (
	"time"

	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type DeliveryStatus string

const (
	DeliveryStatusSent      DeliveryStatus = "sent"
	DeliveryStatusDelivered DeliveryStatus = "delivered"
	DeliveryStatusRead      DeliveryStatus = "read"
)

type Receipt struct {
	ID          uuid.UUID      `json:"id"`
	MessageID   uuid.UUID      `json:"message_id"`
	UserID      uuid.UUID      `json:"user_id"`
	Status      DeliveryStatus `json:"status"`
	DeliveredAt *time.Time     `json:"delivered_at,omitempty"`
	ReadAt      *time.Time     `json:"read_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func NewReceipt(messageID, userID uuid.UUID) *Receipt {
	now := time.Now()
	return &Receipt{
		ID:        uid.New(),
		MessageID: messageID,
		UserID:    userID,
		Status:    DeliveryStatusSent,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (r *Receipt) MarkDelivered() {
	now := time.Now()
	r.Status = DeliveryStatusDelivered
	r.DeliveredAt = &now
	r.UpdatedAt = now
}

func (r *Receipt) MarkRead() {
	now := time.Now()
	if r.DeliveredAt == nil {
		r.DeliveredAt = &now
	}
	r.Status = DeliveryStatusRead
	r.ReadAt = &now
	r.UpdatedAt = now
}

func (r *Receipt) IsDelivered() bool {
	return r.Status == DeliveryStatusDelivered || r.Status == DeliveryStatusRead
}

func (r *Receipt) IsRead() bool {
	return r.Status == DeliveryStatusRead
}
