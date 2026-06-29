package orders_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mmorall/booksonline/internal/orders"

	"github.com/google/uuid"
)

type mockOrderRepository struct {
	onCreate  func(order *orders.Order) error
	onGetByID func(id uuid.UUID) (*orders.Order, error)
	onListAll func() ([]*orders.Order, error)
}

func (m *mockOrderRepository) Create(ctx context.Context, order *orders.Order) error {
	return m.onCreate(order)
}
func (m *mockOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*orders.Order, error) {
	return m.onGetByID(id)
}
func (m *mockOrderRepository) ListAll(ctx context.Context) ([]*orders.Order, error) {
	return m.onListAll()
}

type mockInventoryService struct {
	onGetDetails func(id uuid.UUID) (int64, string, error)
	onReserve    func(id uuid.UUID, q int) error
}

func (m *mockInventoryService) GetProductPriceAndType(ctx context.Context, id uuid.UUID) (int64, string, error) {
	return m.onGetDetails(id)
}
func (m *mockInventoryService) ReserveStock(ctx context.Context, id uuid.UUID, quantity int) error {
	return m.onReserve(id, quantity)
}

func TestOrderService_PlaceOrder_Success(t *testing.T) {
	ctx := context.Background()
	prodID := uuid.New()

	mockInv := &mockInventoryService{
		onGetDetails: func(id uuid.UUID) (int64, string, error) {
			return 2999, "BOOK", nil
		},
		onReserve: func(id uuid.UUID, q int) error {
			return nil // Stock reservation succeeds
		},
	}

	mockRepo := &mockOrderRepository{
		onCreate: func(order *orders.Order) error {
			return nil // Database persistence succeeds
		},
	}

	svc := orders.NewService(mockRepo, mockInv)

	items := map[uuid.UUID]int{prodID: 2}
	order, err := svc.PlaceOrder(ctx, "test@mydomain.com", items)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if order.TotalPrice != 5998 { // 2999 * 2
		t.Errorf("expected total price 5998, got %d", order.TotalPrice)
	}
	if order.Status != orders.StatusCompleted {
		t.Errorf("expected completed status, got %s", order.Status)
	}
}

func TestOrderService_PlaceOrder_OutOfStock(t *testing.T) {
	ctx := context.Background()
	prodID := uuid.New()
	errOutOfStock := errors.New("product out of stock")

	mockInv := &mockInventoryService{
		onGetDetails: func(id uuid.UUID) (int64, string, error) {
			return 9900, "BOOK", nil
		},
		onReserve: func(id uuid.UUID, q int) error {
			return errOutOfStock // Simulate stock depletion
		},
	}
	mockRepo := &mockOrderRepository{} // Unused for this path

	svc := orders.NewService(mockRepo, mockInv)
	items := map[uuid.UUID]int{prodID: 1}

	_, err := svc.PlaceOrder(ctx, "test@mydomain.com", items)

	if !errors.Is(err, errOutOfStock) {
		t.Errorf("expected out of stock error wrapping, got %v", err)
	}
}
