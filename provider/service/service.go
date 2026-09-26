package service

import (
	"context"

	"github.com/ChristianDenniss/go-data-model/provider/entity"
	"github.com/ChristianDenniss/go-data-model/provider/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Record(ctx context.Context, p entity.Provider) error {
	if p.ID == "" {
		return entity.ErrIDRequired
	}
	return s.repo.Upsert(ctx, p)
}

func (s *Service) GetByID(ctx context.Context, id string) (entity.Provider, error) {
	if id == "" {
		return entity.Provider{}, entity.ErrIDRequired
	}
	return s.repo.GetByID(ctx, id)
}
