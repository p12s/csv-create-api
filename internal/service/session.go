package service

import (
	"context"

	"github.com/p12s/csv-create-api/internal/session"
)

// Sessioner - service contract
type Sessioner interface {
	SetExpiredSession(ctx context.Context, sessionToken string, userId int) error
	GetSession(ctx context.Context, sessionToken string) (int, error)
	DeleteSession(ctx context.Context, sessionToken string) error
}

// SessionService - service
type SessionService struct {
	sessionRepo session.Sessioner
}

// NewSessionService - constructor
func NewSessionService(sessionRepo session.Sessioner) *SessionService {
	return &SessionService{sessionRepo: sessionRepo}
}

// SetExpiredSession
func (s *SessionService) SetExpiredSession(ctx context.Context, sessionToken string, userId int) error {
	return s.sessionRepo.SetExpiredSession(ctx, sessionToken, userId)
}

// UpdateProductById
func (s *SessionService) GetSession(ctx context.Context, sessionToken string) (int, error) {
	return s.sessionRepo.GetSession(ctx, sessionToken)
}

// DeleteById
func (s *SessionService) DeleteSession(ctx context.Context, sessionToken string) error {
	return s.sessionRepo.DeleteSession(ctx, sessionToken)
}
