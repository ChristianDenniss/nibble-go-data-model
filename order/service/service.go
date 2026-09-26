package service

import (
	"context"

	"github.com/ChristianDenniss/go-data-model/order/entity"
	"github.com/ChristianDenniss/go-data-model/order/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Record(ctx context.Context, o entity.Order) error {
	if o.ID == "" {
		return entity.ErrIDRequired
	}
	return s.repo.Upsert(ctx, o)
}

func (s *Service) GetByID(ctx context.Context, id string) (entity.Order, error) {
	if id == "" {
		return entity.Order{}, entity.ErrIDRequired
	}
	return s.repo.GetByID(ctx, id)
}
