package service

import (
	"context"

	"github.com/ChristianDenniss/go-data-model/cuisine/entity"
	"github.com/ChristianDenniss/go-data-model/cuisine/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Record(ctx context.Context, c entity.Cuisine) error {
	if c.ID == "" {
		return entity.ErrIDRequired
	}
	return s.repo.Upsert(ctx, c)
}

func (s *Service) GetByID(ctx context.Context, id string) (entity.Cuisine, error) {
	if id == "" {
		return entity.Cuisine{}, entity.ErrIDRequired
	}
	return s.repo.GetByID(ctx, id)
}
