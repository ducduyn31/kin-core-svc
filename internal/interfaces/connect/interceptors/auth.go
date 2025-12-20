package interceptors

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"github.com/danielng/kin-core-svc/internal/application/user"
	"github.com/danielng/kin-core-svc/internal/infrastructure/auth"
)

type ctxKey string

const (
	UserKey     ctxKey = "user"
	UserIDKey   ctxKey = "user_id"
	Auth0SubKey ctxKey = "auth0_sub"
)

type AuthInterceptor struct {
	validator     *auth.Auth0Validator
	userService   *user.Service
	publicMethods map[string]bool
}

func NewAuthInterceptor(validator *auth.Auth0Validator, userService *user.Service) *AuthInterceptor {
	return &AuthInterceptor{
		validator:     validator,
		userService:   userService,
		publicMethods: map[string]bool{},
	}
}

func (i *AuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if i.publicMethods[req.Spec().Procedure] {
			return next(ctx, req)
		}

		newCtx, err := i.authenticate(ctx, req.Header())
		if err != nil {
			return nil, err
		}

		return next(newCtx, req)
	}
}

func (i *AuthInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (i *AuthInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		if i.publicMethods[conn.Spec().Procedure] {
			return next(ctx, conn)
		}

		newCtx, err := i.authenticate(ctx, conn.RequestHeader())
		if err != nil {
			return err
		}

		return next(newCtx, conn)
	}
}

func (i *AuthInterceptor) authenticate(ctx context.Context, headers map[string][]string) (context.Context, error) {
	authHeader := ""
	if values := headers["Authorization"]; len(values) > 0 {
		authHeader = values[0]
	}
	if authHeader == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if token == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	claims, err := i.validator.ValidateToken(ctx, token)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	auth0Sub := claims.GetSub()
	displayName := claims.Name
	if displayName == "" {
		displayName = claims.Email
	}
	if displayName == "" {
		displayName = "User"
	}

	u, err := i.userService.GetOrCreateUser(ctx, auth0Sub, displayName)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	newCtx := context.WithValue(ctx, UserKey, u)
	newCtx = context.WithValue(newCtx, UserIDKey, u.ID)
	newCtx = context.WithValue(newCtx, Auth0SubKey, auth0Sub)

	return newCtx, nil
}
