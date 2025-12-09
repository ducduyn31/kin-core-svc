package availability

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusFree         Status = "free"
	StatusBusy         Status = "busy"
	StatusDoNotDisturb Status = "do_not_disturb"
	StatusSleeping     Status = "sleeping"
	StatusAway         Status = "away"
)

type Availability struct {
	UserID        uuid.UUID  `json:"user_id"`
	Status        Status     `json:"status"`
	StatusMessage *string    `json:"status_message,omitempty"`
	ManualUntil   *time.Time `json:"manual_until,omitempty"` // If manually set, when it expires
	AutoStatus    bool       `json:"auto_status"`            // Whether status is automatically determined
	UpdatedAt     time.Time  `json:"updated_at"`
}

func NewAvailability(userID uuid.UUID) *Availability {
	return &Availability{
		UserID:     userID,
		Status:     StatusFree,
		AutoStatus: true,
		UpdatedAt:  time.Now(),
	}
}

func (a *Availability) SetStatus(status Status, message *string, duration *time.Duration) {
	a.Status = status
	a.StatusMessage = message
	a.AutoStatus = false
	if duration != nil {
		until := time.Now().Add(*duration)
		a.ManualUntil = &until
	} else {
		a.ManualUntil = nil
	}
	a.UpdatedAt = time.Now()
}

func (a *Availability) SetAutoStatus() {
	a.AutoStatus = true
	a.ManualUntil = nil
	a.StatusMessage = nil
	a.UpdatedAt = time.Now()
}

func (a *Availability) IsManualExpired() bool {
	if a.AutoStatus || a.ManualUntil == nil {
		return false
	}
	return time.Now().After(*a.ManualUntil)
}

func IsValidStatus(s Status) bool {
	switch s {
	case StatusFree, StatusBusy, StatusDoNotDisturb, StatusSleeping, StatusAway:
		return true
	default:
		return false
	}
}
