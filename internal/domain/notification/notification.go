package notification

import (
	"time"

	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeMessage        NotificationType = "message"
	NotificationTypeReaction       NotificationType = "reaction"
	NotificationTypeMention        NotificationType = "mention"
	NotificationTypeCircleInvite   NotificationType = "circle_invite"
	NotificationTypeContactRequest NotificationType = "contact_request"
	NotificationTypeCheckIn        NotificationType = "check_in"
	NotificationTypeAvailability   NotificationType = "availability"
)

type Notification struct {
	ID        uuid.UUID         `json:"id"`
	UserID    uuid.UUID         `json:"user_id"`
	Type      NotificationType  `json:"type"`
	Title     string            `json:"title"`
	Body      string            `json:"body"`
	Data      map[string]string `json:"data,omitempty"`
	IsRead    bool              `json:"is_read"`
	IsSent    bool              `json:"is_sent"`
	SentAt    *time.Time        `json:"sent_at,omitempty"`
	ReadAt    *time.Time        `json:"read_at,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

func NewNotification(userID uuid.UUID, notifType NotificationType, title, body string) *Notification {
	return &Notification{
		ID:        uid.New(),
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Body:      body,
		Data:      make(map[string]string),
		IsRead:    false,
		IsSent:    false,
		CreatedAt: time.Now(),
	}
}

func (n *Notification) SetData(key, value string) {
	if n.Data == nil {
		n.Data = make(map[string]string)
	}
	n.Data[key] = value
}

func (n *Notification) MarkSent() {
	now := time.Now()
	n.IsSent = true
	n.SentAt = &now
}

func (n *Notification) MarkRead() {
	now := time.Now()
	n.IsRead = true
	n.ReadAt = &now
}

func NewMessageNotification(userID uuid.UUID, senderName, message string, conversationID uuid.UUID) *Notification {
	notif := NewNotification(userID, NotificationTypeMessage, senderName, message)
	notif.SetData("conversation_id", conversationID.String())
	return notif
}

func NewCircleInviteNotification(userID uuid.UUID, inviterName, circleName string, invitationID uuid.UUID) *Notification {
	notif := NewNotification(
		userID,
		NotificationTypeCircleInvite,
		"Circle Invitation",
		inviterName+" invited you to join "+circleName,
	)
	notif.SetData("invitation_id", invitationID.String())
	return notif
}

func NewContactRequestNotification(userID uuid.UUID, requesterName string, requestID uuid.UUID) *Notification {
	notif := NewNotification(
		userID,
		NotificationTypeContactRequest,
		"Contact Request",
		requesterName+" wants to connect with you",
	)
	notif.SetData("request_id", requestID.String())
	return notif
}

func NewCheckInNotification(userID uuid.UUID, userName, placeName string, checkInID uuid.UUID) *Notification {
	notif := NewNotification(
		userID,
		NotificationTypeCheckIn,
		"Check-in",
		userName+" arrived at "+placeName,
	)
	notif.SetData("check_in_id", checkInID.String())
	return notif
}
