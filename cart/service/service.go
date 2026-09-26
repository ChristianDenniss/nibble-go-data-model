package service

import (
	"context"

	"github.com/ChristianDenniss/go-data-model/cart/entity"
	"github.com/ChristianDenniss/go-data-model/cart/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Record(ctx context.Context, c entity.Cart) error {
	if c.ID == "" {
		return entity.ErrIDRequired
	}
	return s.repo.Upsert(ctx, c)
}

func (s *Service) GetByID(ctx context.Context, id string) (entity.Cart, error) {
	if id == "" {
		return entity.Cart{}, entity.ErrIDRequired
	}
	return s.repo.GetByID(ctx, id)
}
