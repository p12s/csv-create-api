package service

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/p12s/csv-create-api/internal/domain"
	"github.com/p12s/csv-create-api/internal/repository"
)

const (
	TIME_TO_LIVE_HOURS = 24
)

// Authorizer - service contract
type Authorizer interface {
	CreateUser(ctx context.Context, input domain.SignUpInput) error
	GetUserByCredentials(ctx context.Context, email, password string) (string, error)
	ParseToken(ctx context.Context, token string) (int, error)
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
func (a *AuthService) GetUserByCredentials(ctx context.Context, email, password string) (string, error) {
	passwordHash, err := a.generatePasswordHash(password)
	if err != nil {
		return "", fmt.Errorf("get user by creds: %w", err)
	}

	user, err := a.repo.GetByCredentials(ctx, email, passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user with this creds not found")
		}

		return "", fmt.Errorf("user creds wrong: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{ //nolint
		Subject:   strconv.Itoa(user.Id),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(TIME_TO_LIVE_HOURS * time.Hour).Unix(),
	})

	return token.SignedString(a.hmacSecret)
}

func (a *AuthService) generatePasswordHash(password string) (string, error) {
	hash := sha1.New() // #nosec
	if _, err := hash.Write([]byte(password)); err != nil {
		return "", fmt.Errorf("hash write: %w", err)
	}
	return fmt.Sprintf("%x", hash.Sum([]byte(a.salt))), nil
}

func (a *AuthService) ParseToken(ctx context.Context, token string) (int, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return a.hmacSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return id, nil
}
