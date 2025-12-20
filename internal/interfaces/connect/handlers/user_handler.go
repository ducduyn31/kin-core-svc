package handlers

import (
	"context"

	"connectrpc.com/connect"
	kinv1 "github.com/danielng/kin-core-svc/gen/proto/kin/v1"
	"github.com/danielng/kin-core-svc/gen/proto/kin/v1/kinv1connect"
	"github.com/danielng/kin-core-svc/internal/application/user"
	userDomain "github.com/danielng/kin-core-svc/internal/domain/user"
	"github.com/danielng/kin-core-svc/internal/interfaces/connect/converter"
	"github.com/danielng/kin-core-svc/internal/interfaces/connect/interceptors"
	"github.com/google/uuid"
)

type UserHandler struct {
	kinv1connect.UnimplementedUserServiceHandler
	userService *user.Service
}

// NewUserHandler creates a UserHandler backed by the provided user service.
func NewUserHandler(userService *user.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetMe(ctx context.Context, req *connect.Request[kinv1.GetMeRequest]) (*connect.Response[kinv1.GetMeResponse], error) {
	u, ok := ctx.Value(interceptors.UserKey).(*userDomain.User)
	if !ok {
		return nil, connect.NewError(connect.CodeInternal, nil)
	}

	return connect.NewResponse(&kinv1.GetMeResponse{
		User: converter.UserToProto(u),
	}), nil
}

func (h *UserHandler) UpdateProfile(ctx context.Context, req *connect.Request[kinv1.UpdateProfileRequest]) (*connect.Response[kinv1.UpdateProfileResponse], error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, connect.NewError(connect.CodeInternal, nil)
	}

	var displayName string
	if req.Msg.DisplayName != nil {
		displayName = *req.Msg.DisplayName
	}

	u, err := h.userService.UpdateProfile(ctx, user.UpdateProfileCommand{
		UserID:      userID,
		DisplayName: displayName,
		Bio:         req.Msg.Bio,
		Avatar:      req.Msg.Avatar,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(&kinv1.UpdateProfileResponse{
		User: converter.UserToProto(u),
	}), nil
}

func (h *UserHandler) UpdateTimezone(ctx context.Context, req *connect.Request[kinv1.UpdateTimezoneRequest]) (*connect.Response[kinv1.UpdateTimezoneResponse], error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, connect.NewError(connect.CodeInternal, nil)
	}

	if req.Msg.Timezone == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}

	u, err := h.userService.UpdateTimezone(ctx, user.UpdateTimezoneCommand{
		UserID:   userID,
		Timezone: req.Msg.Timezone,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(&kinv1.UpdateTimezoneResponse{
		User: converter.UserToProto(u),
	}), nil
}

func (h *UserHandler) GetPreferences(ctx context.Context, req *connect.Request[kinv1.GetPreferencesRequest]) (*connect.Response[kinv1.GetPreferencesResponse], error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, connect.NewError(connect.CodeInternal, nil)
	}

	prefs, err := h.userService.GetPreferences(ctx, user.GetUserPreferencesQuery{
		UserID: userID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(&kinv1.GetPreferencesResponse{
		Preferences: converter.PreferencesToProto(prefs),
	}), nil
}

func (h *UserHandler) UpdatePreferences(ctx context.Context, req *connect.Request[kinv1.UpdatePreferencesRequest]) (*connect.Response[kinv1.UpdatePreferencesResponse], error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, connect.NewError(connect.CodeInternal, nil)
	}

	cmd := user.UpdatePreferencesCommand{
		UserID:                userID,
		ShowOnlineStatus:      req.Msg.ShowOnlineStatus,
		ShowLastSeen:          req.Msg.ShowLastSeen,
		ShowReadReceipts:      req.Msg.ShowReadReceipts,
		AllowContactDiscovery: req.Msg.AllowContactDiscovery,
		PushNotifications:     req.Msg.PushNotifications,
		EmailNotifications:    req.Msg.EmailNotifications,
		QuietHoursEnabled:     req.Msg.QuietHoursEnabled,
		QuietHoursStart:       req.Msg.QuietHoursStart,
		QuietHoursEnd:         req.Msg.QuietHoursEnd,
	}

	if req.Msg.DefaultPrivacyLevel != nil {
		level := string(converter.PrivacyLevelFromProto(*req.Msg.DefaultPrivacyLevel))
		cmd.DefaultPrivacyLevel = &level
	}

	prefs, err := h.userService.UpdatePreferences(ctx, cmd)
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(&kinv1.UpdatePreferencesResponse{
		Preferences: converter.PreferencesToProto(prefs),
	}), nil
}