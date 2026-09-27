package repository

import (
	"context"

	menuentity "github.com/ChristianDenniss/go-data-model/menu/entity"
	restaurantentity "github.com/ChristianDenniss/go-data-model/restaurant/entity"
	"github.com/ChristianDenniss/go-data-model/storefront/entity"
)

type Repository interface {
	LoadCatalog(ctx context.Context, accountID string) (entity.Catalog, error)
}

// BootstrapRepository loads only the account and browse metadata needed to
// render navigation and search filters. It deliberately excludes the full
// restaurant, menu-item, and offer collections.
type BootstrapRepository interface {
	LoadBootstrap(ctx context.Context, accountID string) (entity.Catalog, error)
}

// RestaurantQuery describes the public catalog search contract. Slugs are
// used at the boundary so callers do not need to know database IDs.
type RestaurantQuery struct {
	Search       string
	CategorySlug string
	CuisineSlug  string
	Page         int
	PageSize     int
	Sort         string
}

type RestaurantPage struct {
	Restaurants []restaurantentity.Restaurant
	Total       int
}

type MenuItemQuery struct {
	Search       string
	RestaurantID string
	CategorySlug string
	CuisineSlug  string
	Page         int
	PageSize     int
}

type MenuItemPage struct {
	Items []menuentity.Item
	Total int
}

// CatalogSearchRepository is optional so existing storefront repository
// implementations remain valid while search is rolled out incrementally.
type CatalogSearchRepository interface {
	SearchRestaurants(ctx context.Context, query RestaurantQuery) (RestaurantPage, error)
	SearchMenuItems(ctx context.Context, query MenuItemQuery) (MenuItemPage, error)
}
