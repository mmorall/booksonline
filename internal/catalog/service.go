package catalog

import (
	"context"

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
