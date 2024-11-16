package session

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

const SESSION_TTL_HOUR = 12

// Sessioner - session work contract
type Sessioner interface {
	SetExpiredSession(ctx context.Context, sessionToken string, userId int) error
	GetSession(ctx context.Context, sessionToken string) (int, error)
	DeleteSession(ctx context.Context, sessionToken string) error
}

// Repository - repo
type Repository struct {
	Sessioner
}

// NewRepository - constructor
func NewRepository(redis *RedisClient) *Repository {
	return &Repository{
		Sessioner: NewSession(redis),
	}
}

// Session
type Session struct {
	client *RedisClient
}

// NewSession - constructor
func NewSession(client *RedisClient) *Session {
	return &Session{client: client}
}

// SetExpiredSession
func (s *Session) SetExpiredSession(ctx context.Context, sessionToken string, userId int) error {
	ttl := SESSION_TTL_HOUR * time.Hour
	err := s.client.client.Set(ctx, sessionToken, userId, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}

	return nil
}

// GetSession
func (s *Session) GetSession(ctx context.Context, sessionToken string) (int, error) {
	val, err := s.client.client.Get(ctx, sessionToken).Result()
	if err != nil {
		fmt.Println("redis get value: %w", err)
	}
	userId, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("redis convert value: %w", err)
	}

	return userId, nil
}

// DeleteSession
func (s *Session) DeleteSession(ctx context.Context, sessionToken string) error {
	err := s.client.client.Del(ctx, sessionToken).Err()
	if err != nil {
		return fmt.Errorf("redis delete error: %w", err)
	}

	return nil
}
