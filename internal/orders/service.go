package orders

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type orderService struct {
	repo      OrderRepository
	inventory InventoryService
}

func NewService(repo OrderRepository, inventory InventoryService) Service {
	return &orderService{
		repo:      repo,
		inventory: inventory,
	}
}

func (s *orderService) PlaceOrder(ctx context.Context, customerEmail string, items map[uuid.UUID]int) (*Order, error) {
	if len(items) == 0 {
		return nil, ErrInvalidOrder
	}

	orderID := uuid.New()
	var orderItems []OrderItem

	for productID, quantity := range items {
		if quantity <= 0 {
			return nil, fmt.Errorf("%w: quantity must be greater than zero", ErrInvalidOrder)
		}

		price, prodType, err := s.inventory.GetProductPriceAndType(ctx, productID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product info for %s: %w", productID, err)
		}

		err = s.inventory.ReserveStock(ctx, productID, quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to reserve stock for %s: %w", productID, err)
		}

		var generatedAsset string

		prefixes := map[string]string{
			"LICENSE": "LIC",
			"VOUCHER": "VCH",
		}

		if prefix, ok := prefixes[prodType]; ok {
			generatedAsset, err = generateDigitalKey(prefix)
			if err != nil {
				return nil, fmt.Errorf("failed to generate digital key for %s: %w", productID, err)
			}
		}

		orderItems = append(orderItems, OrderItem{
			ID:             uuid.New(),
			OrderID:        orderID,
			ProductID:      productID,
			Quantity:       quantity,
			Price:          price,
			GeneratedAsset: generatedAsset,
		})
	}

	order := &Order{
		ID:            orderID,
		CustomerEmail: customerEmail,
		Status:        StatusCompleted, // This would be Pending until payment clears
		CreatedAt:     time.Now().UTC(),
		Items:         orderItems,
	}

	order.CalculateTotal()

	if err := s.repo.Create(ctx, order); err != nil {
		// We could use the Saga pattern or a shared database transaction
		// to rollback the stock reservation.
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return order, nil
}

func (s *orderService) GetOrder(ctx context.Context, id uuid.UUID) (*Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *orderService) ListOrders(ctx context.Context) ([]*Order, error) {
	return s.repo.ListAll(ctx)
}

func generateDigitalKey(prefix string) (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("crypto/rand failed to generate secure key: %w", err)
	}
	return fmt.Sprintf("%s-%s", prefix, hex.EncodeToString(bytes)), nil
}
