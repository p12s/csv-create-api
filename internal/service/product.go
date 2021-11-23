package service

import (
	"context"

	"github.com/p12s/csv-create-api/internal/domain"
)

type Producter interface {
	Create(ctx context.Context, product domain.Product) error
	UpdateById(ctx context.Context, id int, input domain.UpdateProductInput) error
	DeleteById(ctx context.Context, id int) error
	GetAllProducts(ctx context.Context) ([]domain.Product, error)
}

type Service struct {
	repo Producter
}

func NewService(repo Producter) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, product domain.Product) error {
	return s.repo.Create(ctx, product)
}

func (s *Service) UpdateById(ctx context.Context, id int, input domain.UpdateProductInput) error {
	return s.repo.UpdateById(ctx, id, input)
}

func (s *Service) DeleteById(ctx context.Context, id int) error {
	return s.repo.DeleteById(ctx, id)
}

func (s *Service) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	return s.repo.GetAllProducts(ctx)
}
