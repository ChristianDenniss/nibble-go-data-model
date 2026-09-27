package entity

import "time"

// Slots are the ad surfaces a placement can render in. Organic rails are never slots.
const (
	SlotHomeBanner  = "home_banner"
	SlotHomeRail    = "home_rail"
	SlotSearchTop   = "search_top"
	SlotCategoryTop = "category_top"
)

const (
	CampaignDraft  = "draft"
	CampaignActive = "active"
	CampaignPaused = "paused"
	CampaignEnded  = "ended"
)

const (
	PricingCPM  = "cpm"
	PricingCPC  = "cpc"
	PricingFlat = "flat"
)

const (
	AdvertiserActive    = "active"
	AdvertiserSuspended = "suspended"
)

const (
	EventImpression = "impression"
	EventClick      = "click"
)

// Advertiser is a party paying Nibble for placements (a brand, one restaurant, or a channel).
type Advertiser struct {
	ID           string
	Name         string
	BrandID      string
	ContactEmail string
	Status       string
	CreatedAt    time.Time
}

// Campaign is one paid flight: a window, a bid, and budgets. Empty MarketID means all markets.
type Campaign struct {
	ID               string
	AdvertiserID     string
	MarketID         string
	Name             string
	Status           string
	StartsAt         time.Time
	EndsAt           time.Time
	PricingModel     string
	BidCents         int64
	DailyBudgetCents int64
	TotalBudgetCents int64
	Currency         string
}

// Placement is one creative in one slot. Exactly one of LegacyRestaurantID, PlaceID, or BrandID is the target;
// PromotionID is set when the creative features a public deal.
type Placement struct {
	ID                 string
	CampaignID         string
	Slot               string
	Priority           int
	LegacyRestaurantID string
	PlaceID            string
	BrandID            string
	PromotionID        string
	CategoryID         string
	CuisineID          string
	Headline           string
	Body               string
	ImageURL           string
	CallToAction       string
}

// ServedPlacement is a live placement with the campaign fields needed to rank and label it.
type ServedPlacement struct {
	Placement      Placement
	CampaignID     string
	AdvertiserName string
	BidCents       int64
}

// Event is one impression or click; append-only billing source.
type Event struct {
	ID          string
	PlacementID string
	Kind        string
	UserID      string
	Surface     string
	OccurredAt  time.Time
}
