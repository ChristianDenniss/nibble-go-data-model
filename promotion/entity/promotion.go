package entity

import (
	"time"

	money "github.com/ChristianDenniss/go-data-model/money/entity"
)

const (
	KindPercentOff   = "percent_off"
	KindAmountOff    = "amount_off"
	KindFreeDelivery = "free_delivery"
	KindFixedPrice   = "fixed_price"
)

// Promotion is a public deal a channel runs. Empty FulfillmentMode means every mode on that channel.
type Promotion struct {
	ID              string
	ChannelID       string
	Name            string
	Description     string
	Kind            string
	FulfillmentMode string
	Value           money.Money
	ValueBPS        int
	StartsAt        time.Time
	EndsAt          time.Time
}

type Constraint struct {
	ID                 string
	PromotionID        string
	MinSubtotalCents   int64
	Code               string
	MembershipRequired bool
	MaxDiscountCents   int64
}

// Target scopes a promotion. LegacyRestaurantID is an MVP shortcut until storefront reads places.
type Target struct {
	ID                 string
	PromotionID        string
	PlaceID            string
	SourceStoreID      string
	SourceItemID       string
	DishID             string
	BrandID            string
	LegacyRestaurantID string
	Region             string
	Country            string
}

// ActivePromotion is a live promotion with its targets (no targets = channel-wide).
type ActivePromotion struct {
	Promotion Promotion
	Targets   []Target
}

type MembershipProduct struct {
	ID        string
	ChannelID string
	Name      string
	Slug      string
}
