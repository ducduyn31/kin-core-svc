package circle

import (
	"context"
	"log/slog"
	"time"

	"github.com/danielng/kin-core-svc/internal/domain/circle"
	"github.com/google/uuid"
)

type Service struct {
	repo   circle.Repository
	logger *slog.Logger
}

func NewService(repo circle.Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) CreateCircle(ctx context.Context, cmd CreateCircleCommand) (*circle.Circle, error) {
	c := circle.NewCircle(cmd.Name, cmd.Description, cmd.CreatedBy)

	if err := s.repo.Create(ctx, c); err != nil {
		s.logger.Error("failed to create circle", "error", err)
		return nil, err
	}

	member := circle.NewMember(c.ID, cmd.CreatedBy, circle.MemberRoleAdmin)
	if err := s.repo.AddMember(ctx, member); err != nil {
		s.logger.Error("failed to add creator as member", "error", err)
		return nil, err
	}

	pref := circle.NewSharingPreference(c.ID, cmd.CreatedBy)
	if err := s.repo.CreateSharingPreference(ctx, pref); err != nil {
		s.logger.Error("failed to create sharing preferences", "error", err)
	}

	s.logger.Info("circle created", "circle_id", c.ID, "created_by", cmd.CreatedBy)
	return c, nil
}

func (s *Service) GetCircle(ctx context.Context, query GetCircleQuery) (*circle.Circle, error) {
	isMember, err := s.repo.IsMember(ctx, query.CircleID, query.UserID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, circle.ErrNotCircleMember
	}

	return s.repo.GetByID(ctx, query.CircleID)
}

func (s *Service) ListUserCircles(ctx context.Context, query ListUserCirclesQuery) ([]*circle.Circle, int64, error) {
	limit := query.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	circles, err := s.repo.ListByUser(ctx, query.UserID, limit, query.Offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountByUser(ctx, query.UserID)
	if err != nil {
		return nil, 0, err
	}

	return circles, total, nil
}

func (s *Service) UpdateCircle(ctx context.Context, cmd UpdateCircleCommand) (*circle.Circle, error) {
	isAdmin, err := s.repo.IsAdmin(ctx, cmd.CircleID, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, circle.ErrNotCircleAdmin
	}

	c, err := s.repo.GetByID(ctx, cmd.CircleID)
	if err != nil {
		return nil, err
	}

	c.Update(cmd.Name, cmd.Description)

	if err := s.repo.Update(ctx, c); err != nil {
		s.logger.Error("failed to update circle", "error", err, "circle_id", cmd.CircleID)
		return nil, err
	}

	return c, nil
}

func (s *Service) DeleteCircle(ctx context.Context, cmd DeleteCircleCommand) error {
	isAdmin, err := s.repo.IsAdmin(ctx, cmd.CircleID, cmd.UserID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return circle.ErrNotCircleAdmin
	}

	if err := s.repo.Delete(ctx, cmd.CircleID); err != nil {
		s.logger.Error("failed to delete circle", "error", err, "circle_id", cmd.CircleID)
		return err
	}

	s.logger.Info("circle deleted", "circle_id", cmd.CircleID, "deleted_by", cmd.UserID)
	return nil
}

func (s *Service) AddMember(ctx context.Context, cmd AddMemberCommand) (*circle.Member, error) {
	isAdmin, err := s.repo.IsAdmin(ctx, cmd.CircleID, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, circle.ErrNotCircleAdmin
	}

	isMember, err := s.repo.IsMember(ctx, cmd.CircleID, cmd.MemberID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, circle.ErrAlreadyMember
	}

	member := circle.NewMember(cmd.CircleID, cmd.MemberID, cmd.Role)
	if err := s.repo.AddMember(ctx, member); err != nil {
		s.logger.Error("failed to add member", "error", err)
		return nil, err
	}

	pref := circle.NewSharingPreference(cmd.CircleID, cmd.MemberID)
	if err := s.repo.CreateSharingPreference(ctx, pref); err != nil {
		s.logger.Error("failed to create sharing preferences", "error", err)
	}

	s.logger.Info("member added to circle", "circle_id", cmd.CircleID, "member_id", cmd.MemberID)
	return member, nil
}

func (s *Service) RemoveMember(ctx context.Context, cmd RemoveMemberCommand) error {
	isAdmin, err := s.repo.IsAdmin(ctx, cmd.CircleID, cmd.UserID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return circle.ErrNotCircleAdmin
	}

	if cmd.MemberID != cmd.UserID {
		member, err := s.repo.GetMember(ctx, cmd.CircleID, cmd.MemberID)
		if err != nil {
			return err
		}
		if member.IsAdmin() {
			members, err := s.repo.ListMembers(ctx, cmd.CircleID)
			if err != nil {
				return err
			}
			adminCount := 0
			for _, m := range members {
				if m.IsAdmin() {
					adminCount++
				}
			}
			if adminCount <= 1 {
				return circle.ErrCannotRemoveLastAdmin
			}
		}
	}

	if err := s.repo.RemoveMember(ctx, cmd.CircleID, cmd.MemberID); err != nil {
		s.logger.Error("failed to remove member", "error", err)
		return err
	}

	s.logger.Info("member removed from circle", "circle_id", cmd.CircleID, "member_id", cmd.MemberID)
	return nil
}

func (s *Service) ListMembers(ctx context.Context, query ListCircleMembersQuery) ([]*circle.Member, error) {
	isMember, err := s.repo.IsMember(ctx, query.CircleID, query.UserID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, circle.ErrNotCircleMember
	}

	return s.repo.ListMembers(ctx, query.CircleID)
}

func (s *Service) UpdateMemberRole(ctx context.Context, cmd UpdateMemberRoleCommand) (*circle.Member, error) {
	isAdmin, err := s.repo.IsAdmin(ctx, cmd.CircleID, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, circle.ErrNotCircleAdmin
	}

	member, err := s.repo.GetMember(ctx, cmd.CircleID, cmd.MemberID)
	if err != nil {
		return nil, err
	}

	if member.IsAdmin() && cmd.Role != circle.MemberRoleAdmin {
		members, err := s.repo.ListMembers(ctx, cmd.CircleID)
		if err != nil {
			return nil, err
		}
		adminCount := 0
		for _, m := range members {
			if m.IsAdmin() {
				adminCount++
			}
		}
		if adminCount <= 1 {
			return nil, circle.ErrCannotRemoveLastAdmin
		}
	}

	member.SetRole(cmd.Role)

	if err := s.repo.UpdateMember(ctx, member); err != nil {
		s.logger.Error("failed to update member role", "error", err)
		return nil, err
	}

	return member, nil
}

func (s *Service) LeaveCircle(ctx context.Context, cmd LeaveCircleCommand) error {
	member, err := s.repo.GetMember(ctx, cmd.CircleID, cmd.UserID)
	if err != nil {
		return err
	}

	if member.IsAdmin() {
		members, err := s.repo.ListMembers(ctx, cmd.CircleID)
		if err != nil {
			return err
		}
		adminCount := 0
		for _, m := range members {
			if m.IsAdmin() {
				adminCount++
			}
		}
		if adminCount <= 1 {
			return circle.ErrCannotLeaveAsLastAdmin
		}
	}

	if err := s.repo.RemoveMember(ctx, cmd.CircleID, cmd.UserID); err != nil {
		s.logger.Error("failed to leave circle", "error", err)
		return err
	}

	s.logger.Info("user left circle", "circle_id", cmd.CircleID, "user_id", cmd.UserID)
	return nil
}

func (s *Service) GetSharingPreference(ctx context.Context, query GetSharingPreferenceQuery) (*circle.SharingPreference, error) {
	return s.repo.GetSharingPreference(ctx, query.CircleID, query.UserID)
}

func (s *Service) UpdateSharingPreference(ctx context.Context, cmd UpdateSharingPreferenceCommand) (*circle.SharingPreference, error) {
	pref, err := s.repo.GetSharingPreference(ctx, cmd.CircleID, cmd.UserID)
	if err != nil {
		pref = circle.NewSharingPreference(cmd.CircleID, cmd.UserID)
		if err := s.repo.CreateSharingPreference(ctx, pref); err != nil {
			return nil, err
		}
	}

	if cmd.PrivacyLevel != nil {
		pref.PrivacyLevel = *cmd.PrivacyLevel
	}
	if cmd.ShareTimezone != nil {
		pref.ShareTimezone = *cmd.ShareTimezone
	}
	if cmd.ShareAvailability != nil {
		pref.ShareAvailability = *cmd.ShareAvailability
	}
	if cmd.ShareLocation != nil {
		pref.ShareLocation = *cmd.ShareLocation
	}
	if cmd.LocationPrecision != nil {
		pref.LocationPrecision = *cmd.LocationPrecision
	}
	if cmd.ShareActivity != nil {
		pref.ShareActivity = *cmd.ShareActivity
	}
	pref.UpdatedAt = time.Now()

	if err := s.repo.UpdateSharingPreference(ctx, pref); err != nil {
		s.logger.Error("failed to update sharing preference", "error", err)
		return nil, err
	}

	return pref, nil
}

func (s *Service) CreateInvitation(ctx context.Context, cmd CreateInvitationCommand) (*circle.Invitation, error) {
	isAdmin, err := s.repo.IsAdmin(ctx, cmd.CircleID, cmd.InviterID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, circle.ErrNotCircleAdmin
	}

	var expiresAt *time.Time
	if cmd.ExpiresIn != nil {
		t := time.Now().Add(*cmd.ExpiresIn)
		expiresAt = &t
	}

	var inv *circle.Invitation
	if cmd.Type == circle.InvitationTypeDirect && cmd.InviteeID != nil {
		inv = circle.NewDirectInvitation(cmd.CircleID, cmd.InviterID, *cmd.InviteeID, expiresAt)
	} else {
		inv = circle.NewLinkInvitation(cmd.CircleID, cmd.InviterID, cmd.MaxUses, expiresAt)
	}

	if err := s.repo.CreateInvitation(ctx, inv); err != nil {
		s.logger.Error("failed to create invitation", "error", err)
		return nil, err
	}

	s.logger.Info("invitation created", "invitation_id", inv.ID, "circle_id", cmd.CircleID)
	return inv, nil
}

func (s *Service) AcceptInvitation(ctx context.Context, cmd AcceptInvitationCommand) (*circle.Circle, error) {
	inv, err := s.repo.GetInvitationByCode(ctx, cmd.Code)
	if err != nil {
		return nil, err
	}

	if !inv.IsValid() {
		if inv.IsExpired() {
			return nil, circle.ErrInvitationExpired
		}
		return nil, circle.ErrInvitationInvalid
	}

	if inv.Type == circle.InvitationTypeDirect && inv.InviteeID != nil && *inv.InviteeID != cmd.UserID {
		return nil, circle.ErrInvitationInvalid
	}

	isMember, err := s.repo.IsMember(ctx, inv.CircleID, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, circle.ErrAlreadyMember
	}

	member := circle.NewMember(inv.CircleID, cmd.UserID, circle.MemberRoleMember)
	if err := s.repo.AddMember(ctx, member); err != nil {
		s.logger.Error("failed to add member via invitation", "error", err)
		return nil, err
	}

	pref := circle.NewSharingPreference(inv.CircleID, cmd.UserID)
	if err := s.repo.CreateSharingPreference(ctx, pref); err != nil {
		s.logger.Error("failed to create sharing preferences", "error", err)
	}

	inv.Accept()
	if err := s.repo.UpdateInvitation(ctx, inv); err != nil {
		s.logger.Error("failed to update invitation", "error", err)
	}

	s.logger.Info("invitation accepted", "invitation_id", inv.ID, "user_id", cmd.UserID)

	return s.repo.GetByID(ctx, inv.CircleID)
}

func (s *Service) RevokeInvitation(ctx context.Context, cmd RevokeInvitationCommand) error {
	inv, err := s.repo.GetInvitationByID(ctx, cmd.InvitationID)
	if err != nil {
		return err
	}

	isAdmin, err := s.repo.IsAdmin(ctx, inv.CircleID, cmd.UserID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return circle.ErrNotCircleAdmin
	}

	inv.Revoke()
	if err := s.repo.UpdateInvitation(ctx, inv); err != nil {
		s.logger.Error("failed to revoke invitation", "error", err)
		return err
	}

	s.logger.Info("invitation revoked", "invitation_id", cmd.InvitationID)
	return nil
}

func (s *Service) ListUserInvitations(ctx context.Context, query ListUserInvitationsQuery) ([]*circle.Invitation, error) {
	return s.repo.ListUserInvitations(ctx, query.UserID)
}

func (s *Service) GetInvitationByCode(ctx context.Context, query GetInvitationByCodeQuery) (*circle.Invitation, error) {
	return s.repo.GetInvitationByCode(ctx, query.Code)
}

func (s *Service) IsMember(ctx context.Context, circleID, userID uuid.UUID) (bool, error) {
	return s.repo.IsMember(ctx, circleID, userID)
}
