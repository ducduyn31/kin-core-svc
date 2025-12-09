package user

import "github.com/google/uuid"

type CreateUserCommand struct {
	Auth0Sub    string
	DisplayName string
	Email       string
}

type UpdateProfileCommand struct {
	UserID      uuid.UUID
	DisplayName string
	Bio         *string
	Avatar      *string
}

type UpdateTimezoneCommand struct {
	UserID   uuid.UUID
	Timezone string
}

type UpdatePhoneCommand struct {
	UserID      uuid.UUID
	PhoneNumber *string
}

type UpdatePreferencesCommand struct {
	UserID                uuid.UUID
	DefaultPrivacyLevel   *string
	ShowOnlineStatus      *bool
	ShowLastSeen          *bool
	ShowReadReceipts      *bool
	AllowContactDiscovery *bool
	PushNotifications     *bool
	EmailNotifications    *bool
	QuietHoursEnabled     *bool
	QuietHoursStart       *string
	QuietHoursEnd         *string
}
