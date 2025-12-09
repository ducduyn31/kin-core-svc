package user

import "github.com/google/uuid"

type GetUserQuery struct {
	UserID uuid.UUID
}

type GetUserByAuth0SubQuery struct {
	Auth0Sub string
}

type SearchUsersQuery struct {
	Query string
	Limit int
}

type GetUserPreferencesQuery struct {
	UserID uuid.UUID
}
