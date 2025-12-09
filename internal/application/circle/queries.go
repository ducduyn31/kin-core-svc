package circle

import "github.com/google/uuid"

type GetCircleQuery struct {
	CircleID uuid.UUID
	UserID   uuid.UUID // For permission check
}

type ListUserCirclesQuery struct {
	UserID uuid.UUID
	Limit  int
	Offset int
}

type ListCircleMembersQuery struct {
	CircleID uuid.UUID
	UserID   uuid.UUID // For permission check
}

type GetSharingPreferenceQuery struct {
	CircleID uuid.UUID
	UserID   uuid.UUID
}

type ListUserInvitationsQuery struct {
	UserID uuid.UUID
}

type GetInvitationByCodeQuery struct {
	Code string
}
