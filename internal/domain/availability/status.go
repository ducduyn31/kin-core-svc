package availability

import (
	"time"

	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type AutoRule struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Name         string    `json:"name"`
	Condition    Condition `json:"condition"`
	TargetStatus Status    `json:"target_status"`
	Priority     int       `json:"priority"` // Higher priority rules are evaluated first
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ConditionType string

const (
	ConditionTypeTimeRange ConditionType = "time_range" // Based on time of day
	ConditionTypeLocation  ConditionType = "location"   // Based on location
	ConditionTypeCalendar  ConditionType = "calendar"   // Based on calendar events
)

type Condition struct {
	Type ConditionType `json:"type"`

	StartTime *string   `json:"start_time,omitempty"` // HH:MM
	EndTime   *string   `json:"end_time,omitempty"`   // HH:MM
	Weekdays  []Weekday `json:"weekdays,omitempty"`

	PlaceID *uuid.UUID `json:"place_id,omitempty"` // Reference to a user-defined place
}

func NewAutoRule(userID uuid.UUID, name string, condition Condition, targetStatus Status, priority int) *AutoRule {
	now := time.Now()
	return &AutoRule{
		ID:           uid.New(),
		UserID:       userID,
		Name:         name,
		Condition:    condition,
		TargetStatus: targetStatus,
		Priority:     priority,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (r *AutoRule) SetActive(active bool) {
	r.IsActive = active
	r.UpdatedAt = time.Now()
}

func (r *AutoRule) Update(name string, condition Condition, targetStatus Status, priority int) {
	r.Name = name
	r.Condition = condition
	r.TargetStatus = targetStatus
	r.Priority = priority
	r.UpdatedAt = time.Now()
}
