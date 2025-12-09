package circle

import (
	"time"

	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type MemberRole string

const (
	MemberRoleAdmin  MemberRole = "admin"
	MemberRoleMember MemberRole = "member"
)

type Member struct {
	ID        uuid.UUID  `json:"id"`
	CircleID  uuid.UUID  `json:"circle_id"`
	UserID    uuid.UUID  `json:"user_id"`
	Role      MemberRole `json:"role"`
	Nickname  *string    `json:"nickname,omitempty"` // Circle-specific nickname
	JoinedAt  time.Time  `json:"joined_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewMember(circleID, userID uuid.UUID, role MemberRole) *Member {
	now := time.Now()
	return &Member{
		ID:        uid.New(),
		CircleID:  circleID,
		UserID:    userID,
		Role:      role,
		JoinedAt:  now,
		UpdatedAt: now,
	}
}

func (m *Member) IsAdmin() bool {
	return m.Role == MemberRoleAdmin
}

func (m *Member) SetRole(role MemberRole) {
	m.Role = role
	m.UpdatedAt = time.Now()
}

func (m *Member) SetNickname(nickname *string) {
	m.Nickname = nickname
	m.UpdatedAt = time.Now()
}

func IsValidRole(role MemberRole) bool {
	switch role {
	case MemberRoleAdmin, MemberRoleMember:
		return true
	default:
		return false
	}
}
