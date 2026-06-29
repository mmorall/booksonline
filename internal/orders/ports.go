package orders

import (
	"context"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)
	ListAll(ctx context.Context) ([]*Order, error)
}

type InventoryService interface {
	GetProductPriceAndType(ctx context.Context, productID uuid.UUID) (price int64, productType string, err error)
	ReserveStock(ctx context.Context, productID uuid.UUID, quantity int) error
}

type Service interface {
	PlaceOrder(ctx context.Context, customerEmail string, items map[uuid.UUID]int) (*Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*Order, error)
	ListOrders(ctx context.Context) ([]*Order, error)
}
