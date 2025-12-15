package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"connectrpc.com/connect"
	kinv1 "github.com/danielng/kin-core-svc/gen/proto/kin/v1"
	"github.com/danielng/kin-core-svc/gen/proto/kin/v1/kinv1connect"
	"github.com/danielng/kin-core-svc/internal/application/circle"
	"github.com/danielng/kin-core-svc/internal/interfaces/connect/converter"
	"github.com/danielng/kin-core-svc/internal/interfaces/connect/interceptors"
	"github.com/google/uuid"
)

type CircleHandler struct {
	kinv1connect.UnimplementedCircleServiceHandler
	circleService *circle.Service
}

func NewCircleHandler(circleService *circle.Service) *CircleHandler {
	return &CircleHandler{
		circleService: circleService,
	}
}

func (h *CircleHandler) CreateCircle(ctx context.Context, req *connect.Request[kinv1.CreateCircleRequest]) (*connect.Response[kinv1.CreateCircleResponse], error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing user ID in context"))
	}

	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("validation failed: 'name' is required"))
	}

	circ, err := h.circleService.CreateCircle(ctx, circle.CreateCircleCommand{
		Name:        req.Msg.Name,
		Description: req.Msg.Description,
		CreatedBy:   userID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(&kinv1.CreateCircleResponse{
		Circle: converter.CircleToProto(circ),
	}), nil
}

func (h *CircleHandler) ListCircles(ctx context.Context, req *connect.Request[kinv1.ListCirclesRequest]) (*connect.Response[kinv1.ListCirclesResponse], error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing user ID in context"))
	}

	limit := int(req.Msg.Limit)
	if limit <= 0 {
		limit = 20
	}
	offset := int(req.Msg.Offset)
	if offset < 0 {
		offset = 0
	}

	circles, total, err := h.circleService.ListUserCircles(ctx, circle.ListUserCirclesQuery{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(&kinv1.ListCirclesResponse{
		Circles: converter.CirclesToProto(circles),
		Meta: &kinv1.PaginationMeta{
			Page:    int32(offset/limit + 1),
			PerPage: int32(limit),
			Total:   int64(total),
		},
	}), nil
}

func (h *CircleHandler) GetCircle(ctx context.Context, req *connect.Request[kinv1.GetCircleRequest]) (*connect.Response[kinv1.GetCircleResponse], error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing user ID in context"))
	}

	circleID, err := uuid.Parse(req.Msg.CircleId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid UUID for parameter 'circle_id': %w", err))
	}

	circ, err := h.circleService.GetCircle(ctx, circle.GetCircleQuery{
		CircleID: circleID,
		UserID:   userID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(&kinv1.GetCircleResponse{
		Circle: converter.CircleToProto(circ),
	}), nil
}

func (h *CircleHandler) UpdateCircle(ctx context.Context, req *connect.Request[kinv1.UpdateCircleRequest]) (*connect.Response[kinv1.UpdateCircleResponse], error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing user ID in context"))
	}

	circleID, err := uuid.Parse(req.Msg.CircleId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid UUID for parameter 'circle_id': %w", err))
	}

	circ, err := h.circleService.UpdateCircle(ctx, circle.UpdateCircleCommand{
		CircleID:    circleID,
		UserID:      userID,
		Name:        req.Msg.Name,
		Description: req.Msg.Description,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(&kinv1.UpdateCircleResponse{
		Circle: converter.CircleToProto(circ),
	}), nil
}

func (h *CircleHandler) DeleteCircle(ctx context.Context, req *connect.Request[kinv1.DeleteCircleRequest]) (*connect.Response[kinv1.DeleteCircleResponse], error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing user ID in context"))
	}

	circleID, err := uuid.Parse(req.Msg.CircleId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid UUID for parameter 'circle_id': %w", err))
	}

	err = h.circleService.DeleteCircle(ctx, circle.DeleteCircleCommand{
		CircleID: circleID,
		UserID:   userID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(&kinv1.DeleteCircleResponse{}), nil
}

func (h *CircleHandler) LeaveCircle(ctx context.Context, req *connect.Request[kinv1.LeaveCircleRequest]) (*connect.Response[kinv1.LeaveCircleResponse], error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing user ID in context"))
	}

	circleID, err := uuid.Parse(req.Msg.CircleId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid UUID for parameter 'circle_id': %w", err))
	}

	err = h.circleService.LeaveCircle(ctx, circle.LeaveCircleCommand{
		CircleID: circleID,
		UserID:   userID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(&kinv1.LeaveCircleResponse{}), nil
}

func (h *CircleHandler) ListMembers(ctx context.Context, req *connect.Request[kinv1.ListMembersRequest]) (*connect.Response[kinv1.ListMembersResponse], error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing user ID in context"))
	}

	circleID, err := uuid.Parse(req.Msg.CircleId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid UUID for parameter 'circle_id': %w", err))
	}

	members, err := h.circleService.ListMembers(ctx, circle.ListCircleMembersQuery{
		CircleID: circleID,
		UserID:   userID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(&kinv1.ListMembersResponse{
		Members: converter.MembersToProto(members),
	}), nil
}

func (h *CircleHandler) AddMember(ctx context.Context, req *connect.Request[kinv1.AddMemberRequest]) (*connect.Response[kinv1.AddMemberResponse], error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing user ID in context"))
	}

	circleID, err := uuid.Parse(req.Msg.CircleId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid UUID for parameter 'circle_id': %w", err))
	}

	memberID, err := uuid.Parse(req.Msg.UserId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid UUID for parameter 'user_id': %w", err))
	}

	member, err := h.circleService.AddMember(ctx, circle.AddMemberCommand{
		CircleID: circleID,
		UserID:   userID,
		MemberID: memberID,
		Role:     converter.MemberRoleFromProto(req.Msg.Role),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(&kinv1.AddMemberResponse{
		Member: converter.MemberToProto(member),
	}), nil
}

func (h *CircleHandler) RemoveMember(ctx context.Context, req *connect.Request[kinv1.RemoveMemberRequest]) (*connect.Response[kinv1.RemoveMemberResponse], error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing user ID in context"))
	}

	circleID, err := uuid.Parse(req.Msg.CircleId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid UUID for parameter 'circle_id': %w", err))
	}

	memberID, err := uuid.Parse(req.Msg.MemberId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid UUID for parameter 'member_id': %w", err))
	}

	err = h.circleService.RemoveMember(ctx, circle.RemoveMemberCommand{
		CircleID: circleID,
		UserID:   userID,
		MemberID: memberID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(&kinv1.RemoveMemberResponse{}), nil
}

func (h *CircleHandler) GetSharingPreference(ctx context.Context, req *connect.Request[kinv1.GetSharingPreferenceRequest]) (*connect.Response[kinv1.GetSharingPreferenceResponse], error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing user ID in context"))
	}

	circleID, err := uuid.Parse(req.Msg.CircleId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid UUID for parameter 'circle_id': %w", err))
	}

	pref, err := h.circleService.GetSharingPreference(ctx, circle.GetSharingPreferenceQuery{
		CircleID: circleID,
		UserID:   userID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(&kinv1.GetSharingPreferenceResponse{
		SharingPreference: converter.SharingPreferenceToProto(pref),
	}), nil
}

func (h *CircleHandler) UpdateSharingPreference(ctx context.Context, req *connect.Request[kinv1.UpdateSharingPreferenceRequest]) (*connect.Response[kinv1.UpdateSharingPreferenceResponse], error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing user ID in context"))
	}

	circleID, err := uuid.Parse(req.Msg.CircleId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid UUID for parameter 'circle_id': %w", err))
	}

	cmd := circle.UpdateSharingPreferenceCommand{
		CircleID:          circleID,
		UserID:            userID,
		ShareTimezone:     req.Msg.ShareTimezone,
		ShareAvailability: req.Msg.ShareAvailability,
		ShareLocation:     req.Msg.ShareLocation,
		ShareActivity:     req.Msg.ShareActivity,
	}

	if req.Msg.PrivacyLevel != nil {
		level := converter.PrivacyLevelFromProto(*req.Msg.PrivacyLevel)
		cmd.PrivacyLevel = &level
	}

	if req.Msg.LocationPrecision != nil {
		precision := converter.LocationPrecisionFromProto(*req.Msg.LocationPrecision)
		cmd.LocationPrecision = &precision
	}

	pref, err := h.circleService.UpdateSharingPreference(ctx, cmd)
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(&kinv1.UpdateSharingPreferenceResponse{
		SharingPreference: converter.SharingPreferenceToProto(pref),
	}), nil
}

func (h *CircleHandler) CreateInvitation(ctx context.Context, req *connect.Request[kinv1.CreateInvitationRequest]) (*connect.Response[kinv1.CreateInvitationResponse], error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing user ID in context"))
	}

	circleID, err := uuid.Parse(req.Msg.CircleId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid UUID for parameter 'circle_id': %w", err))
	}

	cmd := circle.CreateInvitationCommand{
		CircleID:  circleID,
		InviterID: userID,
		Type:      converter.InvitationTypeFromProto(req.Msg.Type),
	}

	if req.Msg.InviteeId != nil {
		inviteeID, err := uuid.Parse(*req.Msg.InviteeId)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid UUID for parameter 'invitee_id': %w", err))
		}
		cmd.InviteeID = &inviteeID
	}

	if req.Msg.MaxUses != nil {
		maxUses := int(*req.Msg.MaxUses)
		cmd.MaxUses = &maxUses
	}

	if req.Msg.ExpiresInHours != nil {
		duration := time.Duration(*req.Msg.ExpiresInHours) * time.Hour
		cmd.ExpiresIn = &duration
	}

	inv, err := h.circleService.CreateInvitation(ctx, cmd)
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(&kinv1.CreateInvitationResponse{
		Invitation: converter.InvitationToProto(inv),
	}), nil
}

func (h *CircleHandler) AcceptInvitation(ctx context.Context, req *connect.Request[kinv1.AcceptInvitationRequest]) (*connect.Response[kinv1.AcceptInvitationResponse], error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing user ID in context"))
	}

	if req.Msg.Code == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("validation failed: 'code' is required"))
	}

	circ, err := h.circleService.AcceptInvitation(ctx, circle.AcceptInvitationCommand{
		Code:   req.Msg.Code,
		UserID: userID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(&kinv1.AcceptInvitationResponse{
		Circle: converter.CircleToProto(circ),
	}), nil
}
