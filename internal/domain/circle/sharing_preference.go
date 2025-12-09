package circle

import (
	"time"

	"github.com/danielng/kin-core-svc/internal/domain/user"
	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type SharingPreference struct {
	ID                uuid.UUID         `json:"id"`
	CircleID          uuid.UUID         `json:"circle_id"`
	UserID            uuid.UUID         `json:"user_id"`
	PrivacyLevel      user.PrivacyLevel `json:"privacy_level"`
	ShareTimezone     bool              `json:"share_timezone"`
	ShareAvailability bool              `json:"share_availability"`
	ShareLocation     bool              `json:"share_location"`
	LocationPrecision LocationPrecision `json:"location_precision"`
	ShareActivity     bool              `json:"share_activity"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

type LocationPrecision string

const (
	LocationPrecisionCountry      LocationPrecision = "country"
	LocationPrecisionCity         LocationPrecision = "city"
	LocationPrecisionNeighborhood LocationPrecision = "neighborhood"
	LocationPrecisionExact        LocationPrecision = "exact"
)

func NewSharingPreference(circleID, userID uuid.UUID) *SharingPreference {
	now := time.Now()
	return &SharingPreference{
		ID:                uid.New(),
		CircleID:          circleID,
		UserID:            userID,
		PrivacyLevel:      user.PrivacyLevelBasic,
		ShareTimezone:     true,
		ShareAvailability: true,
		ShareLocation:     false,
		LocationPrecision: LocationPrecisionCity,
		ShareActivity:     false,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func (sp *SharingPreference) SetPrivacyLevel(level user.PrivacyLevel) {
	sp.PrivacyLevel = level
	sp.UpdatedAt = time.Now()
}

func (sp *SharingPreference) SetLocationSharing(share bool, precision LocationPrecision) {
	sp.ShareLocation = share
	sp.LocationPrecision = precision
	sp.UpdatedAt = time.Now()
}

func IsValidLocationPrecision(precision LocationPrecision) bool {
	switch precision {
	case LocationPrecisionCountry, LocationPrecisionCity, LocationPrecisionNeighborhood, LocationPrecisionExact:
		return true
	default:
		return false
	}
}
