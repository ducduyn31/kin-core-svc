package user

import (
	"time"

	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	Auth0Sub    string    `json:"auth0_sub"`
	DisplayName string    `json:"display_name"`
	Avatar      *string   `json:"avatar,omitempty"`
	Bio         *string   `json:"bio,omitempty"`
	PhoneNumber *string   `json:"phone_number,omitempty"`
	Timezone    string    `json:"timezone"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewUser(auth0Sub, displayName string) *User {
	now := time.Now()
	return &User{
		ID:          uid.New(),
		Auth0Sub:    auth0Sub,
		DisplayName: displayName,
		Timezone:    "UTC",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (u *User) UpdateProfile(displayName string, bio, avatar *string) {
	if displayName != "" {
		u.DisplayName = displayName
	}
	u.Bio = bio
	u.Avatar = avatar
	u.UpdatedAt = time.Now()
}

func (u *User) SetTimezone(timezone string) {
	u.Timezone = timezone
	u.UpdatedAt = time.Now()
}

func (u *User) SetPhoneNumber(phoneNumber *string) {
	u.PhoneNumber = phoneNumber
	u.UpdatedAt = time.Now()
}

func (u *User) SetAvatar(avatar *string) {
	u.Avatar = avatar
	u.UpdatedAt = time.Now()
}
