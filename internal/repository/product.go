package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/p12s/csv-create-api/internal/domain"
)

// Producter - repository contract
type Producter interface {
	Create(ctx context.Context, product domain.Product) error
	UpdateById(ctx context.Context, id int, input domain.UpdateProductInput) error
	DeleteById(ctx context.Context, id int) error
	GetAllProducts(ctx context.Context) ([]domain.Product, error)
}

// Product - product
type Product struct {
	db *sqlx.DB
}

// NewProduct - constructor
func NewProduct(db *sqlx.DB) *Product {
	return &Product{db: db}
}

// Create - create product
func (p *Product) Create(ctx context.Context, product domain.Product) error {
	query := fmt.Sprintf(`INSERT INTO %s (name, price) values ($1, $2)`, productTable)
	_, err := p.db.Exec(query, product.Name, product.Price)
	if err != nil {
		return fmt.Errorf("create product: %w", err)
	}

	return nil
}

// UpdateById - update product
func (p *Product) UpdateById(ctx context.Context, id int, input domain.UpdateProductInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Price != nil {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argId))
		args = append(args, *input.Price)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d`,
		productTable, setQuery, argId)
	args = append(args, id)

	_, err := p.db.Exec(query, args...)
	return err
}

// DeleteById - delete product
func (p *Product) DeleteById(ctx context.Context, id int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, productTable)
	_, err := p.db.Exec(query, id)
	return err
}

// GetAllProducts - getting all products in csv-file
func (p *Product) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	query := fmt.Sprintf(`SELECT * FROM %s`, productTable)
	if err := p.db.Select(&products, query); err != nil {
		return nil, err
	}

	return products, nil
}
