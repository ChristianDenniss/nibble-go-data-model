package entity

import (
	accountentity "github.com/ChristianDenniss/go-data-model/account/entity"
	cartentity "github.com/ChristianDenniss/go-data-model/cart/entity"
	categoryentity "github.com/ChristianDenniss/go-data-model/category/entity"
	cuisineentity "github.com/ChristianDenniss/go-data-model/cuisine/entity"
	menuentity "github.com/ChristianDenniss/go-data-model/menu/entity"
	offerentity "github.com/ChristianDenniss/go-data-model/offer/entity"
	orderentity "github.com/ChristianDenniss/go-data-model/order/entity"
	promotionentity "github.com/ChristianDenniss/go-data-model/promotion/entity"
	providerentity "github.com/ChristianDenniss/go-data-model/provider/entity"
	restaurantentity "github.com/ChristianDenniss/go-data-model/restaurant/entity"
)

// Catalog is the legacy storefront read model the web app loads in one request.
type Catalog struct {
	Account     accountentity.Account
	Providers   []providerentity.Provider
	Categories  []categoryentity.Category
	Cuisines    []cuisineentity.Cuisine
	Restaurants []restaurantentity.Restaurant
	Items       []menuentity.Item
	Offers      []offerentity.Offer
	Deals       []promotionentity.ActivePromotion
	Cart        cartentity.Cart
	Orders      []orderentity.Order
}
