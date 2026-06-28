package catalog

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrOutOfStock      = errors.New("product is out of stock")
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidProduct  = errors.New("product is inactive or invalid")
)

type ProductType string

const (
	TypeBook    ProductType = "BOOK"
	TypeLicense ProductType = "LICENSE"
	TypeVoucher ProductType = "VOUCHER"
)

type Product struct {
	ID           uuid.UUID
	SKU          string
	Name         string
	PriceInCents int64
	Type         ProductType
	Stock        int
	IsActive     bool
}

func (p *Product) EnsureAvailable() error {
	if !p.IsActive {
		return ErrInvalidProduct
	}
	if p.Type == TypeBook && p.Stock <= 0 {
		return ErrOutOfStock
	}
	return nil
}
