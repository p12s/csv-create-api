package domain

import (
	"time"

	"github.com/go-playground/validator/v10" //nolint
)

var validate *validator.Validate //nolint

// init - init validator
func init() {
	validate = validator.New() //nolint
}

// User
type User struct {
	Id           int        `json:"-" db:"id"`
	Name         string     `json:"name" db:"name" binding:"required"`
	Email        string     `json:"email" db:"email" binding:"required"`
	Password     string     `json:"password" db:"password" binding:"required"`
	RegisteredAt *time.Time `json:"registered_at" db:"registered_at"`
}

// SignUpInput
type SignUpInput struct {
	Name     string `json:"name" validate:"required,gte=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

// Validate
func (i SignUpInput) Validate() error {
	return validate.Struct(i)
}

// SignInInput
type SignInInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

// Validate
func (i SignInInput) Validate() error {
	return validate.Struct(i)
}
