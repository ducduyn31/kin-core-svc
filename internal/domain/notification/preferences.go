package notification

import (
	"time"

	"github.com/google/uuid"
)

type NotificationPreferences struct {
	UserID             uuid.UUID `json:"user_id"`
	PushEnabled        bool      `json:"push_enabled"`
	EmailEnabled       bool      `json:"email_enabled"`
	MessagePush        bool      `json:"message_push"`
	MessageEmail       bool      `json:"message_email"`
	ReactionPush       bool      `json:"reaction_push"`
	CircleInvitePush   bool      `json:"circle_invite_push"`
	CircleInviteEmail  bool      `json:"circle_invite_email"`
	ContactRequestPush bool      `json:"contact_request_push"`
	CheckInPush        bool      `json:"check_in_push"`
	AvailabilityPush   bool      `json:"availability_push"`
	QuietHoursEnabled  bool      `json:"quiet_hours_enabled"`
	QuietHoursStart    *string   `json:"quiet_hours_start,omitempty"` // HH:MM
	QuietHoursEnd      *string   `json:"quiet_hours_end,omitempty"`   // HH:MM
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func NewNotificationPreferences(userID uuid.UUID) *NotificationPreferences {
	now := time.Now()
	return &NotificationPreferences{
		UserID:             userID,
		PushEnabled:        true,
		EmailEnabled:       false,
		MessagePush:        true,
		MessageEmail:       false,
		ReactionPush:       true,
		CircleInvitePush:   true,
		CircleInviteEmail:  true,
		ContactRequestPush: true,
		CheckInPush:        true,
		AvailabilityPush:   true,
		QuietHoursEnabled:  false,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

func (p *NotificationPreferences) IsNotificationEnabled(notifType NotificationType, isPush bool) bool {
	if isPush && !p.PushEnabled {
		return false
	}
	if !isPush && !p.EmailEnabled {
		return false
	}

	switch notifType {
	case NotificationTypeMessage:
		if isPush {
			return p.MessagePush
		}
		return p.MessageEmail
	case NotificationTypeReaction:
		return isPush && p.ReactionPush
	case NotificationTypeCircleInvite:
		if isPush {
			return p.CircleInvitePush
		}
		return p.CircleInviteEmail
	case NotificationTypeContactRequest:
		return isPush && p.ContactRequestPush
	case NotificationTypeCheckIn:
		return isPush && p.CheckInPush
	case NotificationTypeAvailability:
		return isPush && p.AvailabilityPush
	default:
		return true
	}
}

func (p *NotificationPreferences) SetQuietHours(enabled bool, start, end *string) {
	p.QuietHoursEnabled = enabled
	p.QuietHoursStart = start
	p.QuietHoursEnd = end
	p.UpdatedAt = time.Now()
}
