package catalog

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type catalogService struct {
	repo ProductRepository
}

func NewService(repo ProductRepository) Service {
	return &catalogService{
		repo: repo,
	}
}

func (s *catalogService) GetProduct(ctx context.Context, id uuid.UUID) (*Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *catalogService) ListProducts(ctx context.Context) ([]*Product, error) {
	return s.repo.ListAll(ctx)
}

func (s *catalogService) GetProductPriceAndType(ctx context.Context, id uuid.UUID) (int64, string, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return 0, "", err
	}
	return p.PriceInCents, string(p.Type), nil
}

func (s *catalogService) ReserveStock(ctx context.Context, id uuid.UUID, quantity int) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := p.EnsureAvailable(); err != nil {
		return err
	}

	if p.Type == TypeLicense || p.Type == TypeVoucher {
		return nil
	}

	if p.Stock < quantity {
		return ErrOutOfStock
	}
	p.Stock -= quantity

	if err := s.repo.UpdateStock(ctx, p.ID, p.Stock); err != nil {
		return fmt.Errorf("failed to commit stock reservation: %w", err)
	}

	return nil
}
