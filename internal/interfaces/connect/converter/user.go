package converter

import (
	kinv1 "github.com/danielng/kin-core-svc/gen/proto/kin/v1"
	"github.com/danielng/kin-core-svc/internal/domain/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UserToProto(u *user.User) *kinv1.User {
	if u == nil {
		return nil
	}

	pb := &kinv1.User{
		Id:          u.ID.String(),
		Auth0Sub:    u.Auth0Sub,
		DisplayName: u.DisplayName,
		Timezone:    u.Timezone,
		CreatedAt:   timestamppb.New(u.CreatedAt),
		UpdatedAt:   timestamppb.New(u.UpdatedAt),
	}

	if u.Avatar != nil {
		pb.Avatar = u.Avatar
	}
	if u.Bio != nil {
		pb.Bio = u.Bio
	}
	if u.PhoneNumber != nil {
		pb.PhoneNumber = u.PhoneNumber
	}

	return pb
}

func PreferencesToProto(p *user.Preferences) *kinv1.Preferences {
	if p == nil {
		return nil
	}

	pb := &kinv1.Preferences{
		UserId:                p.UserID.String(),
		DefaultPrivacyLevel:   PrivacyLevelToProto(p.DefaultPrivacyLevel),
		ShowOnlineStatus:      p.ShowOnlineStatus,
		ShowLastSeen:          p.ShowLastSeen,
		ShowReadReceipts:      p.ShowReadReceipts,
		AllowContactDiscovery: p.AllowContactDiscovery,
		PushNotifications:     p.PushNotifications,
		EmailNotifications:    p.EmailNotifications,
		QuietHoursEnabled:     p.QuietHoursEnabled,
		CreatedAt:             timestamppb.New(p.CreatedAt),
		UpdatedAt:             timestamppb.New(p.UpdatedAt),
	}

	if p.QuietHoursStart != nil {
		pb.QuietHoursStart = p.QuietHoursStart
	}
	if p.QuietHoursEnd != nil {
		pb.QuietHoursEnd = p.QuietHoursEnd
	}

	return pb
}

func PrivacyLevelToProto(level user.PrivacyLevel) kinv1.PrivacyLevel {
	switch level {
	case user.PrivacyLevelBasic:
		return kinv1.PrivacyLevel_PRIVACY_LEVEL_BASIC
	case user.PrivacyLevelStatus:
		return kinv1.PrivacyLevel_PRIVACY_LEVEL_STATUS
	case user.PrivacyLevelActivity:
		return kinv1.PrivacyLevel_PRIVACY_LEVEL_ACTIVITY
	case user.PrivacyLevelLocation:
		return kinv1.PrivacyLevel_PRIVACY_LEVEL_LOCATION
	default:
		return kinv1.PrivacyLevel_PRIVACY_LEVEL_UNSPECIFIED
	}
}

func PrivacyLevelFromProto(level kinv1.PrivacyLevel) user.PrivacyLevel {
	switch level {
	case kinv1.PrivacyLevel_PRIVACY_LEVEL_BASIC:
		return user.PrivacyLevelBasic
	case kinv1.PrivacyLevel_PRIVACY_LEVEL_STATUS:
		return user.PrivacyLevelStatus
	case kinv1.PrivacyLevel_PRIVACY_LEVEL_ACTIVITY:
		return user.PrivacyLevelActivity
	case kinv1.PrivacyLevel_PRIVACY_LEVEL_LOCATION:
		return user.PrivacyLevelLocation
	default:
		return user.PrivacyLevelBasic
	}
}
