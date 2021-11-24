package service

import (
	"context"

	"github.com/p12s/csv-create-api/internal/domain"
	"github.com/p12s/csv-create-api/internal/repository"
)

// Producter - service contract
type Producter interface {
	CreateProduct(ctx context.Context, product domain.Product) error
	UpdateProductById(ctx context.Context, id int, input domain.UpdateProductInput) error
	DeleteProductById(ctx context.Context, id int) error
	GetAllProducts(ctx context.Context) ([]domain.Product, error)
}

// ProductService - service
type ProductService struct {
	repo repository.Producter
}

// NewProductService - constructor
func NewProductService(repo repository.Producter) *ProductService {
	return &ProductService{repo: repo}
}

// Create
func (s *ProductService) CreateProduct(ctx context.Context, product domain.Product) error {
	return s.repo.Create(ctx, product)
}

// UpdateById
func (s *ProductService) UpdateProductById(ctx context.Context, id int, input domain.UpdateProductInput) error {
	return s.repo.UpdateById(ctx, id, input)
}

// DeleteById
func (s *ProductService) DeleteProductById(ctx context.Context, id int) error {
	return s.repo.DeleteById(ctx, id)
}

// GetAllProducts
func (s *ProductService) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	return s.repo.GetAllProducts(ctx)
}
