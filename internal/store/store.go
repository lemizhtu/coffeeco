package store

import (
	"context"

	"github.com/google/uuid"
	coffeeco "github.com/lemizhtu/coffeego/internal"
)

type Store struct {
	ID              uuid.UUID
	Location        string
	ProductsForSale []coffeeco.Product
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetStoreSpecificDiscount(ctx context.Context, storeID uuid.UUID) (float32, error) {
	discount, err := s.repo.GetStoreDiscount(ctx, storeID)
	if err != nil {
		return 0.0, err
	}

	return discount, nil
}
