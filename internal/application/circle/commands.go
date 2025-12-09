package circle

import (
	"time"

	"github.com/danielng/kin-core-svc/internal/domain/circle"
	"github.com/danielng/kin-core-svc/internal/domain/user"
	"github.com/google/uuid"
)

type CreateCircleCommand struct {
	Name        string
	Description *string
	CreatedBy   uuid.UUID
}

type UpdateCircleCommand struct {
	CircleID    uuid.UUID
	UserID      uuid.UUID // For permission check
	Name        string
	Description *string
}

type DeleteCircleCommand struct {
	CircleID uuid.UUID
	UserID   uuid.UUID // For permission check
}

type AddMemberCommand struct {
	CircleID uuid.UUID
	UserID   uuid.UUID // User adding the member
	MemberID uuid.UUID // User being added
	Role     circle.MemberRole
}

type RemoveMemberCommand struct {
	CircleID uuid.UUID
	UserID   uuid.UUID // User removing the member
	MemberID uuid.UUID // User being removed
}

type UpdateMemberRoleCommand struct {
	CircleID uuid.UUID
	UserID   uuid.UUID // User updating the role
	MemberID uuid.UUID
	Role     circle.MemberRole
}

type LeaveCircleCommand struct {
	CircleID uuid.UUID
	UserID   uuid.UUID
}

type UpdateSharingPreferenceCommand struct {
	CircleID          uuid.UUID
	UserID            uuid.UUID
	PrivacyLevel      *user.PrivacyLevel
	ShareTimezone     *bool
	ShareAvailability *bool
	ShareLocation     *bool
	LocationPrecision *circle.LocationPrecision
	ShareActivity     *bool
}

type CreateInvitationCommand struct {
	CircleID  uuid.UUID
	InviterID uuid.UUID
	InviteeID *uuid.UUID // Nil for link invitations
	Type      circle.InvitationType
	MaxUses   *int
	ExpiresIn *time.Duration
}

type AcceptInvitationCommand struct {
	Code   string
	UserID uuid.UUID
}

type RevokeInvitationCommand struct {
	InvitationID uuid.UUID
	UserID       uuid.UUID // For permission check
}
