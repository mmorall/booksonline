package orders

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidOrder  = errors.New("invalid order data")
	ErrOrderFailed   = errors.New("order processing failed")
	ErrOrderNotFound = errors.New("order not found")
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusCompleted OrderStatus = "COMPLETED"
	StatusFailed    OrderStatus = "FAILED"
)

type Order struct {
	ID            uuid.UUID
	CustomerEmail string
	TotalPrice    int64
	Status        OrderStatus
	CreatedAt     time.Time
	Items         []OrderItem
}

type OrderItem struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	ProductID uuid.UUID
	Quantity  int
	Price     int64

	// For License or Voucher this holds the generated key/code
	GeneratedAsset string
}

func (o *Order) CalculateTotal() {
	var total int64 = 0
	for _, item := range o.Items {
		total += item.Price * int64(item.Quantity)
	}
	o.TotalPrice = total
}
