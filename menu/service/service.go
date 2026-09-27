package service

import (
	"context"

	mediaentity "github.com/ChristianDenniss/go-data-model/media/entity"
	"github.com/ChristianDenniss/go-data-model/menu/entity"
	"github.com/ChristianDenniss/go-data-model/menu/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Record(ctx context.Context, item entity.Item) error {
	if item.ID == "" {
		return entity.ErrIDRequired
	}
	imageURL, err := mediaentity.NormalizeImageURL(item.ImageURL)
	if err != nil {
		return err
	}
	item.ImageURL = imageURL
	return s.repo.Upsert(ctx, item)
}

func (s *Service) GetByID(ctx context.Context, id string) (entity.Item, error) {
	if id == "" {
		return entity.Item{}, entity.ErrIDRequired
	}
	return s.repo.GetByID(ctx, id)
}

// SetImageURL points the item at an external image (empty clears it) and returns the updated item.
func (s *Service) SetImageURL(ctx context.Context, id, rawURL string) (entity.Item, error) {
	if id == "" {
		return entity.Item{}, entity.ErrIDRequired
	}
	imageURL, err := mediaentity.NormalizeImageURL(rawURL)
	if err != nil {
		return entity.Item{}, err
	}
	if err := s.repo.SetImageURL(ctx, id, imageURL); err != nil {
		return entity.Item{}, err
	}
	return s.repo.GetByID(ctx, id)
}
