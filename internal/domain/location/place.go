package location

import (
	"time"

	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type PlaceType string

const (
	PlaceTypeHome   PlaceType = "home"
	PlaceTypeWork   PlaceType = "work"
	PlaceTypeSchool PlaceType = "school"
	PlaceTypeGym    PlaceType = "gym"
	PlaceTypeOther  PlaceType = "other"
)

type Place struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Type      PlaceType `json:"type"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Radius    float64   `json:"radius"` // Detection radius in meters
	Address   *string   `json:"address,omitempty"`
	Icon      *string   `json:"icon,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewPlace(userID uuid.UUID, name string, placeType PlaceType, lat, lng, radius float64) *Place {
	now := time.Now()
	return &Place{
		ID:        uid.New(),
		UserID:    userID,
		Name:      name,
		Type:      placeType,
		Latitude:  lat,
		Longitude: lng,
		Radius:    radius,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (p *Place) Update(name string, placeType PlaceType, lat, lng, radius float64, address *string) {
	p.Name = name
	p.Type = placeType
	p.Latitude = lat
	p.Longitude = lng
	p.Radius = radius
	p.Address = address
	p.UpdatedAt = time.Now()
}

func (p *Place) ContainsLocation(lat, lng float64) bool {
	distance := haversineDistance(p.Latitude, p.Longitude, lat, lng)
	return distance <= p.Radius
}

func IsValidPlaceType(pt PlaceType) bool {
	switch pt {
	case PlaceTypeHome, PlaceTypeWork, PlaceTypeSchool, PlaceTypeGym, PlaceTypeOther:
		return true
	default:
		return false
	}
}
