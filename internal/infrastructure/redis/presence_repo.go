package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/danielng/kin-core-svc/internal/domain/presence"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	presenceKeyPrefix = "presence:"
	activityKeyPrefix = "activity:"
	typingKeyPrefix   = "typing:"
	pushTokenPrefix   = "push_token:"
)

type PresenceRepository struct {
	client *Client
}

func NewPresenceRepository(client *Client) *PresenceRepository {
	return &PresenceRepository{client: client}
}

func (r *PresenceRepository) Set(ctx context.Context, p *presence.Presence, ttl time.Duration) error {
	key := presenceKey(p.UserID)
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal presence: %w", err)
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set presence: %w", err)
	}

	return nil
}

func (r *PresenceRepository) Get(ctx context.Context, userID uuid.UUID) (*presence.Presence, error) {
	key := presenceKey(userID)
	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, presence.ErrPresenceNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get presence: %w", err)
	}

	var p presence.Presence
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("failed to unmarshal presence: %w", err)
	}

	return &p, nil
}

func (r *PresenceRepository) GetMultiple(ctx context.Context, userIDs []uuid.UUID) ([]*presence.Presence, error) {
	if len(userIDs) == 0 {
		return []*presence.Presence{}, nil
	}

	keys := make([]string, len(userIDs))
	for i, id := range userIDs {
		keys[i] = presenceKey(id)
	}

	results, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get multiple presences: %w", err)
	}

	presences := make([]*presence.Presence, 0, len(results))
	for _, result := range results {
		if result == nil {
			continue
		}

		data, ok := result.(string)
		if !ok {
			continue
		}

		var p presence.Presence
		if err := json.Unmarshal([]byte(data), &p); err != nil {
			continue
		}
		presences = append(presences, &p)
	}

	return presences, nil
}

func (r *PresenceRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	key := presenceKey(userID)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete presence: %w", err)
	}
	return nil
}

func (r *PresenceRepository) SetOnline(ctx context.Context, userID uuid.UUID, deviceType presence.DeviceType, deviceID *string, ttl time.Duration) error {
	p := presence.NewPresence(userID)
	p.SetOnline(deviceType, deviceID)
	return r.Set(ctx, p, ttl)
}

func (r *PresenceRepository) SetOffline(ctx context.Context, userID uuid.UUID) error {
	p, err := r.Get(ctx, userID)
	if err != nil {
		p = presence.NewPresence(userID)
	}
	p.SetOffline()

	return r.Set(ctx, p, 24*time.Hour)
}

func (r *PresenceRepository) Heartbeat(ctx context.Context, userID uuid.UUID, ttl time.Duration) error {
	p, err := r.Get(ctx, userID)
	if err != nil {
		return err
	}
	p.Heartbeat()
	return r.Set(ctx, p, ttl)
}

func (r *PresenceRepository) SetActivity(ctx context.Context, a *presence.Activity) error {
	key := activityKey(a.UserID)
	data, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("failed to marshal activity: %w", err)
	}

	ttl := time.Until(a.ExpiresAt)
	if ttl <= 0 {
		return nil
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set activity: %w", err)
	}

	return nil
}

func (r *PresenceRepository) GetActivity(ctx context.Context, userID uuid.UUID) (*presence.Activity, error) {
	key := activityKey(userID)
	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, presence.ErrActivityNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get activity: %w", err)
	}

	var a presence.Activity
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal activity: %w", err)
	}

	return &a, nil
}

func (r *PresenceRepository) ClearActivity(ctx context.Context, userID uuid.UUID) error {
	key := activityKey(userID)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to clear activity: %w", err)
	}
	return nil
}

func (r *PresenceRepository) SetTyping(ctx context.Context, indicator *presence.TypingIndicator) error {
	key := typingKey(indicator.ConversationID)
	ttl := time.Until(indicator.ExpiresAt)
	if ttl <= 0 {
		return nil
	}

	if err := r.client.ZAdd(ctx, key, redis.Z{
		Score:  float64(indicator.ExpiresAt.UnixNano()),
		Member: indicator.UserID.String(),
	}).Err(); err != nil {
		return fmt.Errorf("failed to set typing indicator: %w", err)
	}

	r.client.Expire(ctx, key, ttl+time.Minute)

	return nil
}

func (r *PresenceRepository) GetTypingUsers(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error) {
	key := typingKey(conversationID)
	now := float64(time.Now().UnixNano())

	results, err := r.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: fmt.Sprintf("%f", now),
		Max: "+inf",
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get typing users: %w", err)
	}

	r.client.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%f", now))

	users := make([]uuid.UUID, 0, len(results))
	for _, result := range results {
		userID, err := uuid.Parse(result)
		if err != nil {
			continue
		}
		users = append(users, userID)
	}

	return users, nil
}

func (r *PresenceRepository) ClearTyping(ctx context.Context, userID, conversationID uuid.UUID) error {
	key := typingKey(conversationID)
	if err := r.client.ZRem(ctx, key, userID.String()).Err(); err != nil {
		return fmt.Errorf("failed to clear typing indicator: %w", err)
	}
	return nil
}

func (r *PresenceRepository) SetPushToken(ctx context.Context, userID uuid.UUID, token string) error {
	key := pushTokenKey(userID)
	if err := r.client.Set(ctx, key, token, 0).Err(); err != nil {
		return fmt.Errorf("failed to set push token: %w", err)
	}
	return nil
}

func (r *PresenceRepository) GetPushToken(ctx context.Context, userID uuid.UUID) (string, error) {
	key := pushTokenKey(userID)
	token, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get push token: %w", err)
	}
	return token, nil
}

func (r *PresenceRepository) DeletePushToken(ctx context.Context, userID uuid.UUID) error {
	key := pushTokenKey(userID)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete push token: %w", err)
	}
	return nil
}

func presenceKey(userID uuid.UUID) string {
	return presenceKeyPrefix + userID.String()
}

func activityKey(userID uuid.UUID) string {
	return activityKeyPrefix + userID.String()
}

func typingKey(conversationID uuid.UUID) string {
	return typingKeyPrefix + conversationID.String()
}

func pushTokenKey(userID uuid.UUID) string {
	return pushTokenPrefix + userID.String()
}
