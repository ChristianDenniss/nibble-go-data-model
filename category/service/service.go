package service

import (
	"context"

	"github.com/ChristianDenniss/go-data-model/category/entity"
	"github.com/ChristianDenniss/go-data-model/category/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Record(ctx context.Context, c entity.Category) error {
	if c.ID == "" {
		return entity.ErrIDRequired
	}
	return s.repo.Upsert(ctx, c)
}

func (s *Service) GetByID(ctx context.Context, id string) (entity.Category, error) {
	if id == "" {
		return entity.Category{}, entity.ErrIDRequired
	}
	return s.repo.GetByID(ctx, id)
}
