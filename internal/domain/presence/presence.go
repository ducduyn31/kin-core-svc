package presence

import (
	"time"

	"github.com/google/uuid"
)

type OnlineStatus string

const (
	OnlineStatusOnline  OnlineStatus = "online"
	OnlineStatusOffline OnlineStatus = "offline"
	OnlineStatusAway    OnlineStatus = "away"
)

type DeviceType string

const (
	DeviceTypeMobile  DeviceType = "mobile"
	DeviceTypeDesktop DeviceType = "desktop"
	DeviceTypeWeb     DeviceType = "web"
	DeviceTypeTablet  DeviceType = "tablet"
)

type Presence struct {
	UserID     uuid.UUID    `json:"user_id"`
	Status     OnlineStatus `json:"status"`
	LastSeenAt time.Time    `json:"last_seen_at"`
	DeviceType *DeviceType  `json:"device_type,omitempty"`
	DeviceID   *string      `json:"device_id,omitempty"`
	AppVersion *string      `json:"app_version,omitempty"`
	PushToken  *string      `json:"-"` // Not exposed in API
	UpdatedAt  time.Time    `json:"updated_at"`
}

func NewPresence(userID uuid.UUID) *Presence {
	now := time.Now()
	return &Presence{
		UserID:     userID,
		Status:     OnlineStatusOffline,
		LastSeenAt: now,
		UpdatedAt:  now,
	}
}

func (p *Presence) SetOnline(deviceType DeviceType, deviceID *string) {
	now := time.Now()
	p.Status = OnlineStatusOnline
	p.DeviceType = &deviceType
	p.DeviceID = deviceID
	p.LastSeenAt = now
	p.UpdatedAt = now
}

func (p *Presence) SetOffline() {
	now := time.Now()
	p.Status = OnlineStatusOffline
	p.LastSeenAt = now
	p.UpdatedAt = now
}

func (p *Presence) SetAway() {
	now := time.Now()
	p.Status = OnlineStatusAway
	p.LastSeenAt = now
	p.UpdatedAt = now
}

func (p *Presence) Heartbeat() {
	now := time.Now()
	p.LastSeenAt = now
	p.UpdatedAt = now
}

func (p *Presence) SetPushToken(token *string) {
	p.PushToken = token
	p.UpdatedAt = time.Now()
}

func (p *Presence) SetAppVersion(version *string) {
	p.AppVersion = version
	p.UpdatedAt = time.Now()
}

func (p *Presence) IsOnline() bool {
	return p.Status == OnlineStatusOnline
}

func (p *Presence) IsStale(ttl time.Duration) bool {
	return time.Since(p.LastSeenAt) > ttl
}
