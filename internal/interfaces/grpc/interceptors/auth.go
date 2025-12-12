package interceptors

import (
	"context"
	"strings"

	"github.com/danielng/kin-core-svc/internal/application/user"
	"github.com/danielng/kin-core-svc/internal/infrastructure/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ctxKey string

const (
	UserKey     ctxKey = "user"
	UserIDKey   ctxKey = "user_id"
	Auth0SubKey ctxKey = "auth0_sub"
)

type AuthInterceptor struct {
	validator   *auth.Auth0Validator
	userService *user.Service
	publicMethods map[string]bool
}

func NewAuthInterceptor(validator *auth.Auth0Validator, userService *user.Service) *AuthInterceptor {
	return &AuthInterceptor{
		validator:     validator,
		userService:   userService,
		publicMethods: map[string]bool{},
	}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if i.publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		newCtx, err := i.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		return handler(newCtx, req)
	}
}

func (i *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if i.publicMethods[info.FullMethod] {
			return handler(srv, ss)
		}

		newCtx, err := i.authenticate(ss.Context())
		if err != nil {
			return err
		}

		wrapped := &wrappedStream{
			ServerStream: ss,
			ctx:          newCtx,
		}

		return handler(srv, wrapped)
	}
}

func (i *AuthInterceptor) authenticate(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization header")
	}

	token := strings.TrimPrefix(values[0], "Bearer ")
	claims, err := i.validator.ValidateToken(ctx, token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
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
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	newCtx := context.WithValue(ctx, UserKey, u)
	newCtx = context.WithValue(newCtx, UserIDKey, u.ID)
	newCtx = context.WithValue(newCtx, Auth0SubKey, auth0Sub)

	return newCtx, nil
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}
