package auth

import (
	"slices"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	Scope       string   `json:"scope,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Email       string   `json:"email,omitempty"`
	Name        string   `json:"name,omitempty"`
	Picture     string   `json:"picture,omitempty"`
}

func (c *Claims) HasPermission(permission string) bool {
	return slices.Contains(c.Permissions, permission)
}

func (c *Claims) HasScope(scope string) bool {
	if c.Scope == "" {
		return false
	}
	return c.Scope == scope
}

func (c *Claims) GetSub() string {
	return c.Subject
}
