package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/danielng/kin-core-svc/internal/domain/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (id, auth0_sub, display_name, avatar, bio, phone_number, timezone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Write().Exec(ctx, query, u.ID, u.Auth0Sub, u.DisplayName, u.Avatar, u.Bio, u.PhoneNumber, u.Timezone, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	query := `
		SELECT id, auth0_sub, display_name, avatar, bio, phone_number, timezone, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	return r.scanUser(r.db.Read().QueryRow(ctx, query, id))
}

func (r *UserRepository) GetByAuth0Sub(ctx context.Context, auth0Sub string) (*user.User, error) {
	query := `
		SELECT id, auth0_sub, display_name, avatar, bio, phone_number, timezone, created_at, updated_at
		FROM users
		WHERE auth0_sub = $1
	`
	return r.scanUser(r.db.Read().QueryRow(ctx, query, auth0Sub))
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	query := `
		UPDATE users
		SET display_name = $1, avatar = $2, bio = $3, phone_number = $4, timezone = $5, updated_at = $6
		WHERE id = $7
	`
	_, err := r.db.Write().Exec(ctx, query, u.DisplayName, u.Avatar, u.Bio, u.PhoneNumber, u.Timezone, u.UpdatedAt, u.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*user.User, error) {
	if len(ids) == 0 {
		return []*user.User{}, nil
	}

	query := `
		SELECT id, auth0_sub, display_name, avatar, bio, phone_number, timezone, created_at, updated_at
		FROM users
		WHERE id = ANY($1)
	`
	rows, err := r.db.Read().Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer rows.Close()

	return r.scanUsers(rows)
}

func (r *UserRepository) SearchByDisplayName(ctx context.Context, searchQuery string, limit int) ([]*user.User, error) {
	query := `
		SELECT id, auth0_sub, display_name, avatar, bio, phone_number, timezone, created_at, updated_at
		FROM users
		WHERE display_name ILIKE $1
		LIMIT $2
	`
	rows, err := r.db.Read().Query(ctx, query, "%"+searchQuery+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	return r.scanUsers(rows)
}

func (r *UserRepository) CreatePreferences(ctx context.Context, prefs *user.Preferences) error {
	query := `
		INSERT INTO user_preferences (
			user_id, default_privacy_level, show_online_status, show_last_seen,
			show_read_receipts, allow_contact_discovery, push_notifications, email_notifications,
			quiet_hours_enabled, quiet_hours_start, quiet_hours_end, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.Write().Exec(ctx, query,
		prefs.UserID, prefs.DefaultPrivacyLevel, prefs.ShowOnlineStatus, prefs.ShowLastSeen,
		prefs.ShowReadReceipts, prefs.AllowContactDiscovery, prefs.PushNotifications, prefs.EmailNotifications,
		prefs.QuietHoursEnabled, prefs.QuietHoursStart, prefs.QuietHoursEnd, prefs.CreatedAt, prefs.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user preferences: %w", err)
	}
	return nil
}

func (r *UserRepository) GetPreferences(ctx context.Context, userID uuid.UUID) (*user.Preferences, error) {
	query := `
		SELECT user_id, default_privacy_level, show_online_status, show_last_seen,
			   show_read_receipts, allow_contact_discovery, push_notifications, email_notifications,
			   quiet_hours_enabled, quiet_hours_start, quiet_hours_end, created_at, updated_at
		FROM user_preferences
		WHERE user_id = $1
	`
	row := r.db.Read().QueryRow(ctx, query, userID)

	var prefs user.Preferences
	err := row.Scan(
		&prefs.UserID, &prefs.DefaultPrivacyLevel, &prefs.ShowOnlineStatus, &prefs.ShowLastSeen,
		&prefs.ShowReadReceipts, &prefs.AllowContactDiscovery, &prefs.PushNotifications, &prefs.EmailNotifications,
		&prefs.QuietHoursEnabled, &prefs.QuietHoursStart, &prefs.QuietHoursEnd, &prefs.CreatedAt, &prefs.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, user.ErrPreferencesNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}

	return &prefs, nil
}

func (r *UserRepository) UpdatePreferences(ctx context.Context, prefs *user.Preferences) error {
	query := `
		UPDATE user_preferences
		SET default_privacy_level = $1, show_online_status = $2, show_last_seen = $3,
			show_read_receipts = $4, allow_contact_discovery = $5, push_notifications = $6,
			email_notifications = $7, quiet_hours_enabled = $8, quiet_hours_start = $9,
			quiet_hours_end = $10, updated_at = $11
		WHERE user_id = $12
	`
	_, err := r.db.Write().Exec(ctx, query,
		prefs.DefaultPrivacyLevel, prefs.ShowOnlineStatus, prefs.ShowLastSeen,
		prefs.ShowReadReceipts, prefs.AllowContactDiscovery, prefs.PushNotifications,
		prefs.EmailNotifications, prefs.QuietHoursEnabled, prefs.QuietHoursStart,
		prefs.QuietHoursEnd, prefs.UpdatedAt, prefs.UserID)
	if err != nil {
		return fmt.Errorf("failed to update user preferences: %w", err)
	}
	return nil
}

func (r *UserRepository) scanUser(row pgx.Row) (*user.User, error) {
	var u user.User
	err := row.Scan(
		&u.ID, &u.Auth0Sub, &u.DisplayName, &u.Avatar, &u.Bio,
		&u.PhoneNumber, &u.Timezone, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) scanUsers(rows pgx.Rows) ([]*user.User, error) {
	var users []*user.User
	for rows.Next() {
		var u user.User
		err := rows.Scan(
			&u.ID, &u.Auth0Sub, &u.DisplayName, &u.Avatar, &u.Bio,
			&u.PhoneNumber, &u.Timezone, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &u)
	}
	return users, rows.Err()
}
