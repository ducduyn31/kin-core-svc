package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/danielng/kin-core-svc/internal/domain/circle"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CircleRepository struct {
	db *DB
}

func NewCircleRepository(db *DB) *CircleRepository {
	return &CircleRepository{db: db}
}

func (r *CircleRepository) Create(ctx context.Context, c *circle.Circle) error {
	query := `
		INSERT INTO circles (id, name, description, avatar, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Write().Exec(ctx, query, c.ID, c.Name, c.Description, c.Avatar, c.CreatedBy, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create circle: %w", err)
	}
	return nil
}

func (r *CircleRepository) GetByID(ctx context.Context, id uuid.UUID) (*circle.Circle, error) {
	query := `
		SELECT id, name, description, avatar, created_by, created_at, updated_at
		FROM circles
		WHERE id = $1
	`
	row := r.db.Read().QueryRow(ctx, query, id)

	var c circle.Circle
	err := row.Scan(&c.ID, &c.Name, &c.Description, &c.Avatar, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, circle.ErrCircleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get circle: %w", err)
	}
	return &c, nil
}

func (r *CircleRepository) Update(ctx context.Context, c *circle.Circle) error {
	query := `
		UPDATE circles
		SET name = $1, description = $2, avatar = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.Write().Exec(ctx, query, c.Name, c.Description, c.Avatar, c.UpdatedAt, c.ID)
	if err != nil {
		return fmt.Errorf("failed to update circle: %w", err)
	}
	return nil
}

func (r *CircleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM circles WHERE id = $1`
	_, err := r.db.Write().Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete circle: %w", err)
	}
	return nil
}

func (r *CircleRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*circle.Circle, error) {
	query := `
		SELECT c.id, c.name, c.description, c.avatar, c.created_by, c.created_at, c.updated_at
		FROM circles c
		INNER JOIN circle_members cm ON c.id = cm.circle_id
		WHERE cm.user_id = $1
		ORDER BY c.updated_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Read().Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list circles: %w", err)
	}
	defer rows.Close()

	var circles []*circle.Circle
	for rows.Next() {
		var c circle.Circle
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Avatar, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan circle: %w", err)
		}
		circles = append(circles, &c)
	}
	return circles, rows.Err()
}

func (r *CircleRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM circles c
		INNER JOIN circle_members cm ON c.id = cm.circle_id
		WHERE cm.user_id = $1
	`
	var count int64
	err := r.db.Read().QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count circles: %w", err)
	}
	return count, nil
}

func (r *CircleRepository) AddMember(ctx context.Context, m *circle.Member) error {
	query := `
		INSERT INTO circle_members (id, circle_id, user_id, role, nickname, joined_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Write().Exec(ctx, query, m.ID, m.CircleID, m.UserID, m.Role, m.Nickname, m.JoinedAt, m.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to add circle member: %w", err)
	}
	return nil
}

func (r *CircleRepository) GetMember(ctx context.Context, circleID, userID uuid.UUID) (*circle.Member, error) {
	query := `
		SELECT id, circle_id, user_id, role, nickname, joined_at, updated_at
		FROM circle_members
		WHERE circle_id = $1 AND user_id = $2
	`
	row := r.db.Read().QueryRow(ctx, query, circleID, userID)

	var m circle.Member
	err := row.Scan(&m.ID, &m.CircleID, &m.UserID, &m.Role, &m.Nickname, &m.JoinedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, circle.ErrMemberNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get circle member: %w", err)
	}
	return &m, nil
}

func (r *CircleRepository) GetMemberByID(ctx context.Context, id uuid.UUID) (*circle.Member, error) {
	query := `
		SELECT id, circle_id, user_id, role, nickname, joined_at, updated_at
		FROM circle_members
		WHERE id = $1
	`
	row := r.db.Read().QueryRow(ctx, query, id)

	var m circle.Member
	err := row.Scan(&m.ID, &m.CircleID, &m.UserID, &m.Role, &m.Nickname, &m.JoinedAt, &m.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, circle.ErrMemberNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get circle member: %w", err)
	}
	return &m, nil
}

func (r *CircleRepository) UpdateMember(ctx context.Context, m *circle.Member) error {
	query := `
		UPDATE circle_members
		SET role = $1, nickname = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.Write().Exec(ctx, query, m.Role, m.Nickname, m.UpdatedAt, m.ID)
	if err != nil {
		return fmt.Errorf("failed to update circle member: %w", err)
	}
	return nil
}

func (r *CircleRepository) RemoveMember(ctx context.Context, circleID, userID uuid.UUID) error {
	query := `DELETE FROM circle_members WHERE circle_id = $1 AND user_id = $2`
	_, err := r.db.Write().Exec(ctx, query, circleID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove circle member: %w", err)
	}
	return nil
}

func (r *CircleRepository) ListMembers(ctx context.Context, circleID uuid.UUID) ([]*circle.Member, error) {
	query := `
		SELECT id, circle_id, user_id, role, nickname, joined_at, updated_at
		FROM circle_members
		WHERE circle_id = $1
		ORDER BY joined_at
	`
	rows, err := r.db.Read().Query(ctx, query, circleID)
	if err != nil {
		return nil, fmt.Errorf("failed to list circle members: %w", err)
	}
	defer rows.Close()

	var members []*circle.Member
	for rows.Next() {
		var m circle.Member
		if err := rows.Scan(&m.ID, &m.CircleID, &m.UserID, &m.Role, &m.Nickname, &m.JoinedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan circle member: %w", err)
		}
		members = append(members, &m)
	}
	return members, rows.Err()
}

func (r *CircleRepository) CountMembers(ctx context.Context, circleID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM circle_members WHERE circle_id = $1`
	var count int64
	err := r.db.Read().QueryRow(ctx, query, circleID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count circle members: %w", err)
	}
	return count, nil
}

func (r *CircleRepository) IsMember(ctx context.Context, circleID, userID uuid.UUID) (bool, error) {
	query := `SELECT 1 FROM circle_members WHERE circle_id = $1 AND user_id = $2`
	var exists int
	err := r.db.Read().QueryRow(ctx, query, circleID, userID).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check circle membership: %w", err)
	}
	return true, nil
}

func (r *CircleRepository) IsAdmin(ctx context.Context, circleID, userID uuid.UUID) (bool, error) {
	query := `SELECT 1 FROM circle_members WHERE circle_id = $1 AND user_id = $2 AND role = $3`
	var exists int
	err := r.db.Read().QueryRow(ctx, query, circleID, userID, circle.MemberRoleAdmin).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check circle admin: %w", err)
	}
	return true, nil
}

func (r *CircleRepository) CreateSharingPreference(ctx context.Context, pref *circle.SharingPreference) error {
	query := `
		INSERT INTO circle_sharing_preferences (
			id, circle_id, user_id, privacy_level, share_timezone, share_availability,
			share_location, location_precision, share_activity, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Write().Exec(ctx, query,
		pref.ID, pref.CircleID, pref.UserID, pref.PrivacyLevel, pref.ShareTimezone, pref.ShareAvailability,
		pref.ShareLocation, pref.LocationPrecision, pref.ShareActivity, pref.CreatedAt, pref.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create sharing preference: %w", err)
	}
	return nil
}

func (r *CircleRepository) GetSharingPreference(ctx context.Context, circleID, userID uuid.UUID) (*circle.SharingPreference, error) {
	query := `
		SELECT id, circle_id, user_id, privacy_level, share_timezone, share_availability,
			   share_location, location_precision, share_activity, created_at, updated_at
		FROM circle_sharing_preferences
		WHERE circle_id = $1 AND user_id = $2
	`
	row := r.db.Read().QueryRow(ctx, query, circleID, userID)

	var pref circle.SharingPreference
	err := row.Scan(
		&pref.ID, &pref.CircleID, &pref.UserID, &pref.PrivacyLevel, &pref.ShareTimezone, &pref.ShareAvailability,
		&pref.ShareLocation, &pref.LocationPrecision, &pref.ShareActivity, &pref.CreatedAt, &pref.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, circle.ErrSharingPreferenceNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get sharing preference: %w", err)
	}
	return &pref, nil
}

func (r *CircleRepository) UpdateSharingPreference(ctx context.Context, pref *circle.SharingPreference) error {
	query := `
		UPDATE circle_sharing_preferences
		SET privacy_level = $1, share_timezone = $2, share_availability = $3,
			share_location = $4, location_precision = $5, share_activity = $6, updated_at = $7
		WHERE id = $8
	`
	_, err := r.db.Write().Exec(ctx, query,
		pref.PrivacyLevel, pref.ShareTimezone, pref.ShareAvailability,
		pref.ShareLocation, pref.LocationPrecision, pref.ShareActivity, pref.UpdatedAt, pref.ID)
	if err != nil {
		return fmt.Errorf("failed to update sharing preference: %w", err)
	}
	return nil
}

func (r *CircleRepository) ListSharingPreferences(ctx context.Context, userID uuid.UUID) ([]*circle.SharingPreference, error) {
	query := `
		SELECT id, circle_id, user_id, privacy_level, share_timezone, share_availability,
			   share_location, location_precision, share_activity, created_at, updated_at
		FROM circle_sharing_preferences
		WHERE user_id = $1
	`
	rows, err := r.db.Read().Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list sharing preferences: %w", err)
	}
	defer rows.Close()

	var prefs []*circle.SharingPreference
	for rows.Next() {
		var pref circle.SharingPreference
		if err := rows.Scan(
			&pref.ID, &pref.CircleID, &pref.UserID, &pref.PrivacyLevel, &pref.ShareTimezone, &pref.ShareAvailability,
			&pref.ShareLocation, &pref.LocationPrecision, &pref.ShareActivity, &pref.CreatedAt, &pref.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan sharing preference: %w", err)
		}
		prefs = append(prefs, &pref)
	}
	return prefs, rows.Err()
}

func (r *CircleRepository) CreateInvitation(ctx context.Context, inv *circle.Invitation) error {
	query := `
		INSERT INTO circle_invitations (
			id, circle_id, inviter_id, invitee_id, type, code, status,
			max_uses, use_count, expires_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Write().Exec(ctx, query,
		inv.ID, inv.CircleID, inv.InviterID, inv.InviteeID, inv.Type, inv.Code, inv.Status,
		inv.MaxUses, inv.UseCount, inv.ExpiresAt, inv.CreatedAt, inv.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create invitation: %w", err)
	}
	return nil
}

func (r *CircleRepository) GetInvitationByID(ctx context.Context, id uuid.UUID) (*circle.Invitation, error) {
	query := `
		SELECT id, circle_id, inviter_id, invitee_id, type, code, status,
			   max_uses, use_count, expires_at, created_at, updated_at
		FROM circle_invitations
		WHERE id = $1
	`
	row := r.db.Read().QueryRow(ctx, query, id)
	return r.scanInvitation(row)
}

func (r *CircleRepository) GetInvitationByCode(ctx context.Context, code string) (*circle.Invitation, error) {
	query := `
		SELECT id, circle_id, inviter_id, invitee_id, type, code, status,
			   max_uses, use_count, expires_at, created_at, updated_at
		FROM circle_invitations
		WHERE code = $1
	`
	row := r.db.Read().QueryRow(ctx, query, code)
	return r.scanInvitation(row)
}

func (r *CircleRepository) UpdateInvitation(ctx context.Context, inv *circle.Invitation) error {
	query := `
		UPDATE circle_invitations
		SET status = $1, use_count = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.Write().Exec(ctx, query, inv.Status, inv.UseCount, inv.UpdatedAt, inv.ID)
	if err != nil {
		return fmt.Errorf("failed to update invitation: %w", err)
	}
	return nil
}

func (r *CircleRepository) ListPendingInvitations(ctx context.Context, circleID uuid.UUID) ([]*circle.Invitation, error) {
	query := `
		SELECT id, circle_id, inviter_id, invitee_id, type, code, status,
			   max_uses, use_count, expires_at, created_at, updated_at
		FROM circle_invitations
		WHERE circle_id = $1 AND status = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.Read().Query(ctx, query, circleID, circle.InvitationStatusPending)
	if err != nil {
		return nil, fmt.Errorf("failed to list invitations: %w", err)
	}
	defer rows.Close()
	return r.scanInvitations(rows)
}

func (r *CircleRepository) ListUserInvitations(ctx context.Context, userID uuid.UUID) ([]*circle.Invitation, error) {
	query := `
		SELECT id, circle_id, inviter_id, invitee_id, type, code, status,
			   max_uses, use_count, expires_at, created_at, updated_at
		FROM circle_invitations
		WHERE invitee_id = $1 AND status = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.Read().Query(ctx, query, userID, circle.InvitationStatusPending)
	if err != nil {
		return nil, fmt.Errorf("failed to list user invitations: %w", err)
	}
	defer rows.Close()
	return r.scanInvitations(rows)
}

func (r *CircleRepository) DeleteExpiredInvitations(ctx context.Context) error {
	query := `DELETE FROM circle_invitations WHERE expires_at < $1 AND status = $2`
	_, err := r.db.Write().Exec(ctx, query, time.Now(), circle.InvitationStatusPending)
	if err != nil {
		return fmt.Errorf("failed to delete expired invitations: %w", err)
	}
	return nil
}

func (r *CircleRepository) scanInvitation(row pgx.Row) (*circle.Invitation, error) {
	var inv circle.Invitation
	err := row.Scan(
		&inv.ID, &inv.CircleID, &inv.InviterID, &inv.InviteeID, &inv.Type, &inv.Code, &inv.Status,
		&inv.MaxUses, &inv.UseCount, &inv.ExpiresAt, &inv.CreatedAt, &inv.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, circle.ErrInvitationNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan invitation: %w", err)
	}
	return &inv, nil
}

func (r *CircleRepository) scanInvitations(rows pgx.Rows) ([]*circle.Invitation, error) {
	var invitations []*circle.Invitation
	for rows.Next() {
		var inv circle.Invitation
		if err := rows.Scan(
			&inv.ID, &inv.CircleID, &inv.InviterID, &inv.InviteeID, &inv.Type, &inv.Code, &inv.Status,
			&inv.MaxUses, &inv.UseCount, &inv.ExpiresAt, &inv.CreatedAt, &inv.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan invitation: %w", err)
		}
		invitations = append(invitations, &inv)
	}
	return invitations, rows.Err()
}
