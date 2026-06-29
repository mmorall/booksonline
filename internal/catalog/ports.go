package catalog

import (
	"context"

	"github.com/google/uuid"
)

type ProductRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	ListAll(ctx context.Context) ([]*Product, error)
	UpdateStock(ctx context.Context, id uuid.UUID, newStock int) error
}

type Service interface {
	GetProduct(ctx context.Context, id uuid.UUID) (*Product, error)
	ListProducts(ctx context.Context) ([]*Product, error)
	GetProductPriceAndType(ctx context.Context, id uuid.UUID) (int64, string, error)
	ReserveStock(ctx context.Context, id uuid.UUID, quantity int) error
}
