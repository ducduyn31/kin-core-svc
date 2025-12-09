package availability

import (
	"time"

	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type Window struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Weekday   Weekday   `json:"weekday"`
	StartTime string    `json:"start_time"` // HH:MM format in user's timezone
	EndTime   string    `json:"end_time"`   // HH:MM format in user's timezone
	Status    Status    `json:"status"`     // Status during this window
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewWindow(userID uuid.UUID, name string, weekday Weekday, startTime, endTime string, status Status) *Window {
	now := time.Now()
	return &Window{
		ID:        uid.New(),
		UserID:    userID,
		Name:      name,
		Weekday:   weekday,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    status,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (w *Window) Update(name string, weekday Weekday, startTime, endTime string, status Status) {
	w.Name = name
	w.Weekday = weekday
	w.StartTime = startTime
	w.EndTime = endTime
	w.Status = status
	w.UpdatedAt = time.Now()
}

func (w *Window) SetActive(active bool) {
	w.IsActive = active
	w.UpdatedAt = time.Now()
}

func IsValidWeekday(d Weekday) bool {
	return d >= Sunday && d <= Saturday
}

func WeekdayFromTime(t time.Time) Weekday {
	return Weekday(t.Weekday())
}
