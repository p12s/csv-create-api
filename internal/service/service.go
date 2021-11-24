package service

import (
	"github.com/p12s/csv-create-api/internal/repository"
)

// Service
type Service struct {
	Producter
	Authorizer
}

// NewService - constructor
func NewService(repos *repository.Repository) *Service {
	return &Service{
		Producter:  NewProductService(repos.Producter),
		Authorizer: NewAuthService(repos.Authorizer),
	}
}
