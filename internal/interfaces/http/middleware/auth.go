package middleware

import (
	"strings"

	"github.com/danielng/kin-core-svc/internal/application/user"
	"github.com/danielng/kin-core-svc/internal/infrastructure/auth"
	"github.com/danielng/kin-core-svc/pkg/ctxkey"
	"github.com/danielng/kin-core-svc/pkg/response"
	"github.com/gin-gonic/gin"
)

func Auth(validator *auth.Auth0Validator, userService *user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := validator.ValidateToken(c.Request.Context(), token)
		if err != nil {
			response.Unauthorized(c, "invalid token")
			c.Abort()
			return
		}

		auth0Sub := claims.GetSub()
		displayName := claims.Name
		if displayName == "" {
			displayName = claims.Email
		}
		if displayName == "" {
			displayName = "User"
		}

		u, err := userService.GetOrCreateUser(c.Request.Context(), auth0Sub, displayName)
		if err != nil {
			response.InternalError(c)
			c.Abort()
			return
		}

		c.Set(ctxkey.User, u)
		c.Set(ctxkey.UserID, u.ID)
		c.Set(ctxkey.Auth0Sub, auth0Sub)

		c.Next()
	}
}

func OptionalAuth(validator *auth.Auth0Validator, userService *user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := validator.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.Next()
			return
		}

		auth0Sub := claims.GetSub()
		u, err := userService.GetOrCreateUser(c.Request.Context(), auth0Sub, claims.Name)
		if err != nil {
			c.Next()
			return
		}

		c.Set(ctxkey.User, u)
		c.Set(ctxkey.UserID, u.ID)
		c.Set(ctxkey.Auth0Sub, auth0Sub)

		c.Next()
	}
}
