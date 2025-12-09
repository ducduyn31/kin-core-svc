package user

import (
	"time"

	"github.com/google/uuid"
)

type PrivacyLevel string

const (
	PrivacyLevelBasic    PrivacyLevel = "basic"    // Online/offline, timezone
	PrivacyLevelStatus   PrivacyLevel = "status"   // At home, At work, Commuting, Sleeping
	PrivacyLevelActivity PrivacyLevel = "activity" // Current activity, estimated free time
	PrivacyLevelLocation PrivacyLevel = "location" // Real-time location sharing
)

type Preferences struct {
	UserID                uuid.UUID    `json:"user_id"`
	DefaultPrivacyLevel   PrivacyLevel `json:"default_privacy_level"`
	ShowOnlineStatus      bool         `json:"show_online_status"`
	ShowLastSeen          bool         `json:"show_last_seen"`
	ShowReadReceipts      bool         `json:"show_read_receipts"`
	AllowContactDiscovery bool         `json:"allow_contact_discovery"`
	PushNotifications     bool         `json:"push_notifications"`
	EmailNotifications    bool         `json:"email_notifications"`
	QuietHoursEnabled     bool         `json:"quiet_hours_enabled"`
	QuietHoursStart       *string      `json:"quiet_hours_start,omitempty"` // HH:MM format
	QuietHoursEnd         *string      `json:"quiet_hours_end,omitempty"`   // HH:MM format
	CreatedAt             time.Time    `json:"created_at"`
	UpdatedAt             time.Time    `json:"updated_at"`
}

func NewPreferences(userID uuid.UUID) *Preferences {
	now := time.Now()
	return &Preferences{
		UserID:                userID,
		DefaultPrivacyLevel:   PrivacyLevelBasic,
		ShowOnlineStatus:      true,
		ShowLastSeen:          true,
		ShowReadReceipts:      true,
		AllowContactDiscovery: true,
		PushNotifications:     true,
		EmailNotifications:    false,
		QuietHoursEnabled:     false,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
}

func (p *Preferences) SetPrivacyLevel(level PrivacyLevel) {
	p.DefaultPrivacyLevel = level
	p.UpdatedAt = time.Now()
}

func (p *Preferences) SetQuietHours(enabled bool, start, end *string) {
	p.QuietHoursEnabled = enabled
	p.QuietHoursStart = start
	p.QuietHoursEnd = end
	p.UpdatedAt = time.Now()
}

func IsValidPrivacyLevel(level PrivacyLevel) bool {
	switch level {
	case PrivacyLevelBasic, PrivacyLevelStatus, PrivacyLevelActivity, PrivacyLevelLocation:
		return true
	default:
		return false
	}
}
