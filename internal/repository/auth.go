package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/p12s/csv-create-api/internal/domain"
)

// Producter - auth contract
type Authorizer interface {
	CreateUser(ctx context.Context, input domain.SignUpInput) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
}

// Auth - Auth
type Auth struct {
	db *sqlx.DB
}

// NewAuth - constructor
func NewAuth(db *sqlx.DB) *Auth {
	return &Auth{db: db}
}

// CreateUser - create user
func (a *Auth) CreateUser(ctx context.Context, input domain.SignUpInput) error {
	query := fmt.Sprintf(`INSERT INTO %s (name, email, password)
		values ($1, $2, $3)`, usersTable)
	_, err := a.db.Exec(query, input.Name, input.Email, input.Password)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

// GetByCredentials - get with credentials
func (a *Auth) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf(`SELECT * FROM %s WHERE email=$1 AND password=$2`, usersTable)
	err := a.db.Get(&user, query, email, password)
	if err != nil {
		return user, fmt.Errorf("get user by creds: %w", err)
	}

	return user, nil
}
