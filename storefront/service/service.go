package service

import (
	"context"
	"errors"

	accountentity "github.com/ChristianDenniss/go-data-model/account/entity"
	"github.com/ChristianDenniss/go-data-model/storefront/entity"
	"github.com/ChristianDenniss/go-data-model/storefront/repository"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) LoadCatalog(ctx context.Context, accountID string) (entity.Catalog, error) {
	if accountID == "" {
		return entity.Catalog{}, entity.ErrAccountIDRequired
	}
	catalog, err := s.repo.LoadCatalog(ctx, accountID)
	if err != nil {
		if errors.Is(err, accountentity.ErrNotFound) {
			return entity.Catalog{}, entity.ErrNotFound
		}
		return entity.Catalog{}, err
	}
	return catalog, nil
}

func (s *Service) LoadBootstrap(ctx context.Context, accountID string) (entity.Catalog, error) {
	if accountID == "" {
		return entity.Catalog{}, entity.ErrAccountIDRequired
	}
	bootstrap, ok := s.repo.(repository.BootstrapRepository)
	if !ok {
		return entity.Catalog{}, errors.New("catalog bootstrap unavailable")
	}
	catalog, err := bootstrap.LoadBootstrap(ctx, accountID)
	if err != nil {
		if errors.Is(err, accountentity.ErrNotFound) {
			return entity.Catalog{}, entity.ErrNotFound
		}
		return entity.Catalog{}, err
	}
	return catalog, nil
}

func (s *Service) SearchRestaurants(ctx context.Context, query repository.RestaurantQuery) (repository.RestaurantPage, error) {
	search, ok := s.repo.(repository.CatalogSearchRepository)
	if !ok {
		return repository.RestaurantPage{}, errors.New("catalog restaurant search unavailable")
	}
	return search.SearchRestaurants(ctx, query)
}

func (s *Service) SearchMenuItems(ctx context.Context, query repository.MenuItemQuery) (repository.MenuItemPage, error) {
	search, ok := s.repo.(repository.CatalogSearchRepository)
	if !ok {
		return repository.MenuItemPage{}, errors.New("catalog menu search unavailable")
	}
	return search.SearchMenuItems(ctx, query)
}
