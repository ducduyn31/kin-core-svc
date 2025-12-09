package circle

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type InvitationType string

const (
	InvitationTypeDirect InvitationType = "direct" // Direct invite to a specific user
	InvitationTypeLink   InvitationType = "link"   // Invite link anyone can use
)

type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusExpired  InvitationStatus = "expired"
	InvitationStatusRevoked  InvitationStatus = "revoked"
)

type Invitation struct {
	ID        uuid.UUID        `json:"id"`
	CircleID  uuid.UUID        `json:"circle_id"`
	InviterID uuid.UUID        `json:"inviter_id"`
	InviteeID *uuid.UUID       `json:"invitee_id,omitempty"` // Null for link invitations
	Type      InvitationType   `json:"type"`
	Code      string           `json:"code"` // Unique invitation code
	Status    InvitationStatus `json:"status"`
	MaxUses   *int             `json:"max_uses,omitempty"`
	UseCount  int              `json:"use_count"`
	ExpiresAt *time.Time       `json:"expires_at,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func NewDirectInvitation(circleID, inviterID, inviteeID uuid.UUID, expiresAt *time.Time) *Invitation {
	now := time.Now()
	return &Invitation{
		ID:        uid.New(),
		CircleID:  circleID,
		InviterID: inviterID,
		InviteeID: &inviteeID,
		Type:      InvitationTypeDirect,
		Code:      generateInviteCode(),
		Status:    InvitationStatusPending,
		UseCount:  0,
		ExpiresAt: expiresAt,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewLinkInvitation(circleID, inviterID uuid.UUID, maxUses *int, expiresAt *time.Time) *Invitation {
	now := time.Now()
	return &Invitation{
		ID:        uid.New(),
		CircleID:  circleID,
		InviterID: inviterID,
		InviteeID: nil,
		Type:      InvitationTypeLink,
		Code:      generateInviteCode(),
		Status:    InvitationStatusPending,
		MaxUses:   maxUses,
		UseCount:  0,
		ExpiresAt: expiresAt,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (i *Invitation) IsValid() bool {
	if i.Status != InvitationStatusPending {
		return false
	}
	if i.ExpiresAt != nil && time.Now().After(*i.ExpiresAt) {
		return false
	}
	if i.MaxUses != nil && i.UseCount >= *i.MaxUses {
		return false
	}
	return true
}

func (i *Invitation) IsExpired() bool {
	if i.Status == InvitationStatusExpired {
		return true
	}
	if i.ExpiresAt != nil && time.Now().After(*i.ExpiresAt) {
		return true
	}
	return false
}

func (i *Invitation) Accept() {
	i.UseCount++
	if i.Type == InvitationTypeDirect || (i.MaxUses != nil && i.UseCount >= *i.MaxUses) {
		i.Status = InvitationStatusAccepted
	}
	i.UpdatedAt = time.Now()
}

func (i *Invitation) Revoke() {
	i.Status = InvitationStatusRevoked
	i.UpdatedAt = time.Now()
}

func generateInviteCode() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:22]
}
