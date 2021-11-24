package service

import (
	"github.com/p12s/csv-create-api/internal/repository"
	sessionRepository "github.com/p12s/csv-create-api/internal/session"
)

// Service
type Service struct {
	Producter
	Authorizer
	Sessioner
}

// NewService - constructor
func NewService(repos *repository.Repository, sessionRepo *sessionRepository.Repository) *Service {
	return &Service{
		Producter:  NewProductService(repos.Producter),
		Authorizer: NewAuthService(repos.Authorizer),
		Sessioner:  NewSessionService(sessionRepo.Sessioner),
	}
}
