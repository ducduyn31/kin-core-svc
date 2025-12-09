package user

import (
	"context"
	"log/slog"

	"github.com/danielng/kin-core-svc/internal/domain/user"
	"github.com/google/uuid"
)

type Service struct {
	repo   user.Repository
	logger *slog.Logger
}

func NewService(repo user.Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) CreateUser(ctx context.Context, cmd CreateUserCommand) (*user.User, error) {
	existingUser, err := s.repo.GetByAuth0Sub(ctx, cmd.Auth0Sub)
	if err == nil {
		return existingUser, nil
	}

	u := user.NewUser(cmd.Auth0Sub, cmd.DisplayName)

	if err := s.repo.Create(ctx, u); err != nil {
		s.logger.Error("failed to create user", "error", err, "auth0_sub", cmd.Auth0Sub)
		return nil, err
	}

	prefs := user.NewPreferences(u.ID)
	if err := s.repo.CreatePreferences(ctx, prefs); err != nil {
		s.logger.Error("failed to create user preferences", "error", err, "user_id", u.ID)
	}

	s.logger.Info("user created", "user_id", u.ID, "auth0_sub", cmd.Auth0Sub)
	return u, nil
}

func (s *Service) GetUser(ctx context.Context, query GetUserQuery) (*user.User, error) {
	return s.repo.GetByID(ctx, query.UserID)
}

func (s *Service) GetUserByAuth0Sub(ctx context.Context, query GetUserByAuth0SubQuery) (*user.User, error) {
	return s.repo.GetByAuth0Sub(ctx, query.Auth0Sub)
}

func (s *Service) GetOrCreateUser(ctx context.Context, auth0Sub, displayName string) (*user.User, error) {
	u, err := s.repo.GetByAuth0Sub(ctx, auth0Sub)
	if err == nil {
		return u, nil
	}

	return s.CreateUser(ctx, CreateUserCommand{
		Auth0Sub:    auth0Sub,
		DisplayName: displayName,
	})
}

func (s *Service) UpdateProfile(ctx context.Context, cmd UpdateProfileCommand) (*user.User, error) {
	u, err := s.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	u.UpdateProfile(cmd.DisplayName, cmd.Bio, cmd.Avatar)

	if err := s.repo.Update(ctx, u); err != nil {
		s.logger.Error("failed to update user profile", "error", err, "user_id", cmd.UserID)
		return nil, err
	}

	return u, nil
}

func (s *Service) UpdateTimezone(ctx context.Context, cmd UpdateTimezoneCommand) (*user.User, error) {
	u, err := s.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	u.SetTimezone(cmd.Timezone)

	if err := s.repo.Update(ctx, u); err != nil {
		s.logger.Error("failed to update user timezone", "error", err, "user_id", cmd.UserID)
		return nil, err
	}

	return u, nil
}

func (s *Service) UpdatePhone(ctx context.Context, cmd UpdatePhoneCommand) (*user.User, error) {
	u, err := s.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	u.SetPhoneNumber(cmd.PhoneNumber)

	if err := s.repo.Update(ctx, u); err != nil {
		s.logger.Error("failed to update user phone", "error", err, "user_id", cmd.UserID)
		return nil, err
	}

	return u, nil
}

func (s *Service) GetPreferences(ctx context.Context, query GetUserPreferencesQuery) (*user.Preferences, error) {
	return s.repo.GetPreferences(ctx, query.UserID)
}

func (s *Service) UpdatePreferences(ctx context.Context, cmd UpdatePreferencesCommand) (*user.Preferences, error) {
	prefs, err := s.repo.GetPreferences(ctx, cmd.UserID)
	if err != nil {
		prefs = user.NewPreferences(cmd.UserID)
		if err := s.repo.CreatePreferences(ctx, prefs); err != nil {
			return nil, err
		}
	}

	if cmd.DefaultPrivacyLevel != nil && user.IsValidPrivacyLevel(user.PrivacyLevel(*cmd.DefaultPrivacyLevel)) {
		prefs.DefaultPrivacyLevel = user.PrivacyLevel(*cmd.DefaultPrivacyLevel)
	}
	if cmd.ShowOnlineStatus != nil {
		prefs.ShowOnlineStatus = *cmd.ShowOnlineStatus
	}
	if cmd.ShowLastSeen != nil {
		prefs.ShowLastSeen = *cmd.ShowLastSeen
	}
	if cmd.ShowReadReceipts != nil {
		prefs.ShowReadReceipts = *cmd.ShowReadReceipts
	}
	if cmd.AllowContactDiscovery != nil {
		prefs.AllowContactDiscovery = *cmd.AllowContactDiscovery
	}
	if cmd.PushNotifications != nil {
		prefs.PushNotifications = *cmd.PushNotifications
	}
	if cmd.EmailNotifications != nil {
		prefs.EmailNotifications = *cmd.EmailNotifications
	}
	if cmd.QuietHoursEnabled != nil {
		prefs.QuietHoursEnabled = *cmd.QuietHoursEnabled
	}
	if cmd.QuietHoursStart != nil {
		prefs.QuietHoursStart = cmd.QuietHoursStart
	}
	if cmd.QuietHoursEnd != nil {
		prefs.QuietHoursEnd = cmd.QuietHoursEnd
	}

	if err := s.repo.UpdatePreferences(ctx, prefs); err != nil {
		s.logger.Error("failed to update preferences", "error", err, "user_id", cmd.UserID)
		return nil, err
	}

	return prefs, nil
}

func (s *Service) SearchUsers(ctx context.Context, query SearchUsersQuery) ([]*user.User, error) {
	limit := query.Limit
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	return s.repo.SearchByDisplayName(ctx, query.Query, limit)
}

func (s *Service) GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]*user.User, error) {
	return s.repo.FindByIDs(ctx, ids)
}
