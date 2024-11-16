package service

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/p12s/csv-create-api/internal/domain"
	"github.com/p12s/csv-create-api/internal/repository"
)

const (
	TIME_TO_LIVE_HOURS = 24
)

// Authorizer - service contract
type Authorizer interface {
	CreateUser(ctx context.Context, input domain.SignUpInput) error
	GetUserByCredentials(ctx context.Context, email, password string) (int, error)
}

// AuthService - service
type AuthService struct {
	repo       repository.Authorizer
	salt       string
	hmacSecret []byte
}

// NewAuthService - constructor
func NewAuthService(repo repository.Authorizer) *AuthService {
	return &AuthService{
		repo:       repo,
		salt:       os.Getenv("SALT"),
		hmacSecret: []byte(os.Getenv("HMAC_SECRET")),
	}
}

// Create
func (a *AuthService) CreateUser(ctx context.Context, input domain.SignUpInput) error {
	passwordHash, err := a.generatePasswordHash(input.Password)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	input.Password = passwordHash
	return a.repo.CreateUser(ctx, input)
}

// GetUserByCredentials
func (a *AuthService) GetUserByCredentials(ctx context.Context, email, password string) (int, error) {
	passwordHash, err := a.generatePasswordHash(password)
	if err != nil {
		return 0, fmt.Errorf("get user by creds: %w", err)
	}

	user, err := a.repo.GetByCredentials(ctx, email, passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user with this creds not found")
		}

		return 0, fmt.Errorf("user creds wrong: %w", err)
	}

	return user.Id, nil
}

// generatePasswordHash
func (a *AuthService) generatePasswordHash(password string) (string, error) {
	hash := sha1.New() // #nosec
	if _, err := hash.Write([]byte(password)); err != nil {
		return "", fmt.Errorf("hash write: %w", err)
	}
	return fmt.Sprintf("%x", hash.Sum([]byte(a.salt))), nil
}
