package circle

import (
	"time"

	"github.com/danielng/kin-core-svc/pkg/uid"
	"github.com/google/uuid"
)

type Circle struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Avatar      *string   `json:"avatar,omitempty"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewCircle(name string, description *string, createdBy uuid.UUID) *Circle {
	now := time.Now()
	return &Circle{
		ID:          uid.New(),
		Name:        name,
		Description: description,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (c *Circle) Update(name string, description *string) {
	if name != "" {
		c.Name = name
	}
	c.Description = description
	c.UpdatedAt = time.Now()
}

func (c *Circle) SetAvatar(avatar *string) {
	c.Avatar = avatar
	c.UpdatedAt = time.Now()
}
