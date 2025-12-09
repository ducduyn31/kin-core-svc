package location

import (
	"time"

	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type Location struct {
	UserID       uuid.UUID  `json:"user_id"`
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
	Accuracy     *float64   `json:"accuracy,omitempty"` // Meters
	Altitude     *float64   `json:"altitude,omitempty"` // Meters
	Speed        *float64   `json:"speed,omitempty"`    // Meters per second
	Heading      *float64   `json:"heading,omitempty"`  // Degrees from north
	PlaceID      *uuid.UUID `json:"place_id,omitempty"` // If at a known place
	Country      *string    `json:"country,omitempty"`
	City         *string    `json:"city,omitempty"`
	Neighborhood *string    `json:"neighborhood,omitempty"`
	Address      *string    `json:"address,omitempty"`
	IsMoving     bool       `json:"is_moving"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func NewLocation(userID uuid.UUID, lat, lng float64) *Location {
	return &Location{
		UserID:    userID,
		Latitude:  lat,
		Longitude: lng,
		IsMoving:  false,
		UpdatedAt: time.Now(),
	}
}

func (l *Location) Update(lat, lng float64, accuracy, altitude, speed, heading *float64) {
	l.Latitude = lat
	l.Longitude = lng
	l.Accuracy = accuracy
	l.Altitude = altitude
	l.Speed = speed
	l.Heading = heading
	l.UpdatedAt = time.Now()

	if speed != nil && *speed > 1.0 { // More than 1 m/s
		l.IsMoving = true
	} else {
		l.IsMoving = false
	}
}

func (l *Location) SetPlace(placeID *uuid.UUID) {
	l.PlaceID = placeID
	l.UpdatedAt = time.Now()
}

func (l *Location) SetGeocodedInfo(country, city, neighborhood, address *string) {
	l.Country = country
	l.City = city
	l.Neighborhood = neighborhood
	l.Address = address
	l.UpdatedAt = time.Now()
}

func (l *Location) DistanceTo(other *Location) float64 {
	return haversineDistance(l.Latitude, l.Longitude, other.Latitude, other.Longitude)
}

type LocationHistory struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Latitude  float64    `json:"latitude"`
	Longitude float64    `json:"longitude"`
	Accuracy  *float64   `json:"accuracy,omitempty"`
	PlaceID   *uuid.UUID `json:"place_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

func NewLocationHistory(userID uuid.UUID, lat, lng float64, accuracy *float64, placeID *uuid.UUID) *LocationHistory {
	return &LocationHistory{
		ID:        uid.New(),
		UserID:    userID,
		Latitude:  lat,
		Longitude: lng,
		Accuracy:  accuracy,
		PlaceID:   placeID,
		CreatedAt: time.Now(),
	}
}
