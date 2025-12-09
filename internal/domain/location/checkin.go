package location

import (
	"time"

	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type CheckInType string

const (
	CheckInTypeArrival   CheckInType = "arrival"
	CheckInTypeDeparture CheckInType = "departure"
)

type CheckIn struct {
	ID        uuid.UUID   `json:"id"`
	UserID    uuid.UUID   `json:"user_id"`
	PlaceID   uuid.UUID   `json:"place_id"`
	Type      CheckInType `json:"type"`
	Latitude  float64     `json:"latitude"`
	Longitude float64     `json:"longitude"`
	Note      *string     `json:"note,omitempty"`
	AutoCheck bool        `json:"auto_check"` // Whether it was automatic or manual
	CreatedAt time.Time   `json:"created_at"`
}

func NewCheckIn(userID, placeID uuid.UUID, checkInType CheckInType, lat, lng float64, autoCheck bool) *CheckIn {
	return &CheckIn{
		ID:        uid.New(),
		UserID:    userID,
		PlaceID:   placeID,
		Type:      checkInType,
		Latitude:  lat,
		Longitude: lng,
		AutoCheck: autoCheck,
		CreatedAt: time.Now(),
	}
}

func (c *CheckIn) SetNote(note *string) {
	c.Note = note
}

func (c *CheckIn) IsArrival() bool {
	return c.Type == CheckInTypeArrival
}

func (c *CheckIn) IsDeparture() bool {
	return c.Type == CheckInTypeDeparture
}
