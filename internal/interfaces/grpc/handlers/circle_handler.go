package handlers

import (
	"context"
	"time"

	kinv1 "github.com/danielng/kin-core-svc/gen/proto/kin/v1"
	"github.com/danielng/kin-core-svc/internal/application/circle"
	circleDomain "github.com/danielng/kin-core-svc/internal/domain/circle"
	"github.com/danielng/kin-core-svc/internal/domain/user"
	"github.com/danielng/kin-core-svc/internal/interfaces/grpc/converter"
	"github.com/danielng/kin-core-svc/internal/interfaces/grpc/interceptors"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CircleHandler struct {
	kinv1.UnimplementedCircleServiceServer
	circleService *circle.Service
}

func NewCircleHandler(circleService *circle.Service) *CircleHandler {
	return &CircleHandler{
		circleService: circleService,
	}
}

func (h *CircleHandler) CreateCircle(ctx context.Context, req *kinv1.CreateCircleRequest) (*kinv1.CreateCircleResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	circ, err := h.circleService.CreateCircle(ctx, circle.CreateCircleCommand{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   userID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &kinv1.CreateCircleResponse{
		Circle: converter.CircleToProto(circ),
	}, nil
}

func (h *CircleHandler) ListCircles(ctx context.Context, req *kinv1.ListCirclesRequest) (*kinv1.ListCirclesResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 20
	}
	offset := int(req.Offset)
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

	return &kinv1.ListCirclesResponse{
		Circles: converter.CirclesToProto(circles),
		Meta: &kinv1.PaginationMeta{
			Page:    int32(offset/limit + 1),
			PerPage: int32(limit),
			Total:   int64(total),
		},
	}, nil
}

func (h *CircleHandler) GetCircle(ctx context.Context, req *kinv1.GetCircleRequest) (*kinv1.GetCircleResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	circleID, err := uuid.Parse(req.CircleId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid circle ID")
	}

	circ, err := h.circleService.GetCircle(ctx, circle.GetCircleQuery{
		CircleID: circleID,
		UserID:   userID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &kinv1.GetCircleResponse{
		Circle: converter.CircleToProto(circ),
	}, nil
}

func (h *CircleHandler) UpdateCircle(ctx context.Context, req *kinv1.UpdateCircleRequest) (*kinv1.UpdateCircleResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	circleID, err := uuid.Parse(req.CircleId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid circle ID")
	}

	var name string
	if req.Name != nil {
		name = *req.Name
	}

	circ, err := h.circleService.UpdateCircle(ctx, circle.UpdateCircleCommand{
		CircleID:    circleID,
		UserID:      userID,
		Name:        name,
		Description: req.Description,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &kinv1.UpdateCircleResponse{
		Circle: converter.CircleToProto(circ),
	}, nil
}

func (h *CircleHandler) DeleteCircle(ctx context.Context, req *kinv1.DeleteCircleRequest) (*kinv1.DeleteCircleResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	circleID, err := uuid.Parse(req.CircleId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid circle ID")
	}

	err = h.circleService.DeleteCircle(ctx, circle.DeleteCircleCommand{
		CircleID: circleID,
		UserID:   userID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &kinv1.DeleteCircleResponse{}, nil
}

func (h *CircleHandler) LeaveCircle(ctx context.Context, req *kinv1.LeaveCircleRequest) (*kinv1.LeaveCircleResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	circleID, err := uuid.Parse(req.CircleId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid circle ID")
	}

	err = h.circleService.LeaveCircle(ctx, circle.LeaveCircleCommand{
		CircleID: circleID,
		UserID:   userID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &kinv1.LeaveCircleResponse{}, nil
}

func (h *CircleHandler) ListMembers(ctx context.Context, req *kinv1.ListMembersRequest) (*kinv1.ListMembersResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	circleID, err := uuid.Parse(req.CircleId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid circle ID")
	}

	members, err := h.circleService.ListMembers(ctx, circle.ListCircleMembersQuery{
		CircleID: circleID,
		UserID:   userID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &kinv1.ListMembersResponse{
		Members: converter.MembersToProto(members),
	}, nil
}

func (h *CircleHandler) AddMember(ctx context.Context, req *kinv1.AddMemberRequest) (*kinv1.AddMemberResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	circleID, err := uuid.Parse(req.CircleId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid circle ID")
	}

	memberID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	member, err := h.circleService.AddMember(ctx, circle.AddMemberCommand{
		CircleID: circleID,
		UserID:   userID,
		MemberID: memberID,
		Role:     converter.MemberRoleFromProto(req.Role),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &kinv1.AddMemberResponse{
		Member: converter.MemberToProto(member),
	}, nil
}

func (h *CircleHandler) RemoveMember(ctx context.Context, req *kinv1.RemoveMemberRequest) (*kinv1.RemoveMemberResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	circleID, err := uuid.Parse(req.CircleId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid circle ID")
	}

	memberID, err := uuid.Parse(req.MemberId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid member ID")
	}

	err = h.circleService.RemoveMember(ctx, circle.RemoveMemberCommand{
		CircleID: circleID,
		UserID:   userID,
		MemberID: memberID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &kinv1.RemoveMemberResponse{}, nil
}

func (h *CircleHandler) GetSharingPreference(ctx context.Context, req *kinv1.GetSharingPreferenceRequest) (*kinv1.GetSharingPreferenceResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	circleID, err := uuid.Parse(req.CircleId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid circle ID")
	}

	pref, err := h.circleService.GetSharingPreference(ctx, circle.GetSharingPreferenceQuery{
		CircleID: circleID,
		UserID:   userID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &kinv1.GetSharingPreferenceResponse{
		SharingPreference: converter.SharingPreferenceToProto(pref),
	}, nil
}

func (h *CircleHandler) UpdateSharingPreference(ctx context.Context, req *kinv1.UpdateSharingPreferenceRequest) (*kinv1.UpdateSharingPreferenceResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	circleID, err := uuid.Parse(req.CircleId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid circle ID")
	}

	cmd := circle.UpdateSharingPreferenceCommand{
		CircleID:          circleID,
		UserID:            userID,
		ShareTimezone:     req.ShareTimezone,
		ShareAvailability: req.ShareAvailability,
		ShareLocation:     req.ShareLocation,
		ShareActivity:     req.ShareActivity,
	}

	if req.PrivacyLevel != nil {
		level := converter.PrivacyLevelFromProto(*req.PrivacyLevel)
		cmd.PrivacyLevel = &level
	}

	if req.LocationPrecision != nil {
		precision := converter.LocationPrecisionFromProto(*req.LocationPrecision)
		cmd.LocationPrecision = &precision
	}

	pref, err := h.circleService.UpdateSharingPreference(ctx, cmd)
	if err != nil {
		return nil, mapError(err)
	}

	return &kinv1.UpdateSharingPreferenceResponse{
		SharingPreference: converter.SharingPreferenceToProto(pref),
	}, nil
}

func (h *CircleHandler) CreateInvitation(ctx context.Context, req *kinv1.CreateInvitationRequest) (*kinv1.CreateInvitationResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	circleID, err := uuid.Parse(req.CircleId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid circle ID")
	}

	cmd := circle.CreateInvitationCommand{
		CircleID:  circleID,
		InviterID: userID,
		Type:      converter.InvitationTypeFromProto(req.Type),
	}

	if req.InviteeId != nil {
		inviteeID, err := uuid.Parse(*req.InviteeId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid invitee ID")
		}
		cmd.InviteeID = &inviteeID
	}

	if req.MaxUses != nil {
		maxUses := int(*req.MaxUses)
		cmd.MaxUses = &maxUses
	}

	if req.ExpiresInHours != nil {
		duration := time.Duration(*req.ExpiresInHours) * time.Hour
		cmd.ExpiresIn = &duration
	}

	inv, err := h.circleService.CreateInvitation(ctx, cmd)
	if err != nil {
		return nil, mapError(err)
	}

	return &kinv1.CreateInvitationResponse{
		Invitation: converter.InvitationToProto(inv),
	}, nil
}

func (h *CircleHandler) AcceptInvitation(ctx context.Context, req *kinv1.AcceptInvitationRequest) (*kinv1.AcceptInvitationResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "user ID not found in context")
	}

	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "invitation code is required")
	}

	circ, err := h.circleService.AcceptInvitation(ctx, circle.AcceptInvitationCommand{
		Code:   req.Code,
		UserID: userID,
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &kinv1.AcceptInvitationResponse{
		Circle: converter.CircleToProto(circ),
	}, nil
}

var _ user.PrivacyLevel
var _ circleDomain.LocationPrecision
