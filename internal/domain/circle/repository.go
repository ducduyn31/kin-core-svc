package circle

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, circle *Circle) error
	GetByID(ctx context.Context, id uuid.UUID) (*Circle, error)
	Update(ctx context.Context, circle *Circle) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Circle, error)
	CountByUser(ctx context.Context, userID uuid.UUID) (int64, error)

	AddMember(ctx context.Context, member *Member) error
	GetMember(ctx context.Context, circleID, userID uuid.UUID) (*Member, error)
	GetMemberByID(ctx context.Context, id uuid.UUID) (*Member, error)
	UpdateMember(ctx context.Context, member *Member) error
	RemoveMember(ctx context.Context, circleID, userID uuid.UUID) error
	ListMembers(ctx context.Context, circleID uuid.UUID) ([]*Member, error)
	CountMembers(ctx context.Context, circleID uuid.UUID) (int64, error)
	IsMember(ctx context.Context, circleID, userID uuid.UUID) (bool, error)
	IsAdmin(ctx context.Context, circleID, userID uuid.UUID) (bool, error)

	CreateSharingPreference(ctx context.Context, pref *SharingPreference) error
	GetSharingPreference(ctx context.Context, circleID, userID uuid.UUID) (*SharingPreference, error)
	UpdateSharingPreference(ctx context.Context, pref *SharingPreference) error
	ListSharingPreferences(ctx context.Context, userID uuid.UUID) ([]*SharingPreference, error)

	CreateInvitation(ctx context.Context, invitation *Invitation) error
	GetInvitationByID(ctx context.Context, id uuid.UUID) (*Invitation, error)
	GetInvitationByCode(ctx context.Context, code string) (*Invitation, error)
	UpdateInvitation(ctx context.Context, invitation *Invitation) error
	ListPendingInvitations(ctx context.Context, circleID uuid.UUID) ([]*Invitation, error)
	ListUserInvitations(ctx context.Context, userID uuid.UUID) ([]*Invitation, error)
	DeleteExpiredInvitations(ctx context.Context) error
}
