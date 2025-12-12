package handlers

import (
	"context"

	kinv1 "github.com/danielng/kin-core-svc/gen/proto/kin/v1"
	"github.com/danielng/kin-core-svc/internal/application/user"
	userDomain "github.com/danielng/kin-core-svc/internal/domain/user"
	"github.com/danielng/kin-core-svc/internal/interfaces/grpc/converter"
	"github.com/danielng/kin-core-svc/internal/interfaces/grpc/interceptors"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	kinv1.UnimplementedUserServiceServer
	userService *user.Service
}

func NewUserHandler(userService *user.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetMe(ctx context.Context, _ *kinv1.GetMeRequest) (*kinv1.GetMeResponse, error) {
	u, ok := ctx.Value(interceptors.UserKey).(*userDomain.User)
	if !ok {
		return nil, status.Error(codes.Internal, "user not found in context")
	}

	return &kinv1.GetMeResponse{
		User: converter.UserToProto(u),
	}, nil
}

func (h *UserHandler) UpdateProfile(ctx context.Context, req *kinv1.UpdateProfileRequest) (*kinv1.UpdateProfileResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	var displayName string
	if req.DisplayName != nil {
		displayName = *req.DisplayName
	}

	u, err := h.userService.UpdateProfile(ctx, user.UpdateProfileCommand{
		UserID:      userID,
		DisplayName: displayName,
		Bio:         req.Bio,
		Avatar:      req.Avatar,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &kinv1.UpdateProfileResponse{
		User: converter.UserToProto(u),
	}, nil
}

func (h *UserHandler) UpdateTimezone(ctx context.Context, req *kinv1.UpdateTimezoneRequest) (*kinv1.UpdateTimezoneResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	if req.Timezone == "" {
		return nil, status.Error(codes.InvalidArgument, "timezone is required")
	}

	u, err := h.userService.UpdateTimezone(ctx, user.UpdateTimezoneCommand{
		UserID:   userID,
		Timezone: req.Timezone,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &kinv1.UpdateTimezoneResponse{
		User: converter.UserToProto(u),
	}, nil
}

func (h *UserHandler) GetPreferences(ctx context.Context, _ *kinv1.GetPreferencesRequest) (*kinv1.GetPreferencesResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	prefs, err := h.userService.GetPreferences(ctx, user.GetUserPreferencesQuery{
		UserID: userID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &kinv1.GetPreferencesResponse{
		Preferences: converter.PreferencesToProto(prefs),
	}, nil
}

func (h *UserHandler) UpdatePreferences(ctx context.Context, req *kinv1.UpdatePreferencesRequest) (*kinv1.UpdatePreferencesResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	cmd := user.UpdatePreferencesCommand{
		UserID:            userID,
		ShowOnlineStatus:  req.ShowOnlineStatus,
		ShowLastSeen:      req.ShowLastSeen,
		ShowReadReceipts:  req.ShowReadReceipts,
		AllowContactDiscovery: req.AllowContactDiscovery,
		PushNotifications:     req.PushNotifications,
		EmailNotifications:    req.EmailNotifications,
		QuietHoursEnabled:     req.QuietHoursEnabled,
		QuietHoursStart:       req.QuietHoursStart,
		QuietHoursEnd:         req.QuietHoursEnd,
	}

	if req.DefaultPrivacyLevel != nil {
		level := string(converter.PrivacyLevelFromProto(*req.DefaultPrivacyLevel))
		cmd.DefaultPrivacyLevel = &level
	}

	prefs, err := h.userService.UpdatePreferences(ctx, cmd)
	if err != nil {
		return nil, mapError(err)
	}

	return &kinv1.UpdatePreferencesResponse{
		Preferences: converter.PreferencesToProto(prefs),
	}, nil
}
