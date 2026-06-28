package catalog_test

import (
	"errors"
	"testing"

	"github.com/mmorall/booksonline/internal/catalog"

	"github.com/google/uuid"
)

func TestProduct_EnsureAvailable(t *testing.T) {
	tests := []struct {
		name          string
		product       catalog.Product
		expectedError error
	}{
		{
			name: "Active book in stock is available",
			product: catalog.Product{
				ID: uuid.New(), Type: catalog.TypeBook, Stock: 5, IsActive: true,
			},
			expectedError: nil,
		},
		{
			name: "Active book out of stock returns error",
			product: catalog.Product{
				ID: uuid.New(), Type: catalog.TypeBook, Stock: 0, IsActive: true,
			},
			expectedError: catalog.ErrOutOfStock,
		},
		{
			name: "Active license with zero stock is available (infinite supply)",
			product: catalog.Product{
				ID: uuid.New(), Type: catalog.TypeLicense, Stock: 0, IsActive: true,
			},
			expectedError: nil,
		},
		{
			name: "Inactive product returns error",
			product: catalog.Product{
				ID: uuid.New(), Type: catalog.TypeBook, Stock: 10, IsActive: false,
			},
			expectedError: catalog.ErrInvalidProduct,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.EnsureAvailable()
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}
