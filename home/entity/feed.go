package entity

import "time"

const (
	SectionSponsored   = "sponsored"
	SectionDeals       = "deals"
	SectionPopular     = "popular"
	SectionRecommended = "recommended"
)

const (
	BannerSponsored = "sponsored"
	BannerDeal      = "deal"
)

// SponsoredMark labels paid content. Anything carrying one must render a visible "Sponsored" tag (D28).
type SponsoredMark struct {
	PlacementID    string
	CampaignID     string
	AdvertiserName string
	Label          string
}

// DealBadge is a provider's public promotion on a restaurant; not paid placement.
type DealBadge struct {
	PromotionID     string
	ChannelID       string
	Label           string
	FulfillmentMode string
}

type Banner struct {
	ID           string
	Kind         string
	Headline     string
	Body         string
	ImageURL     string
	CallToAction string
	RestaurantID string
	Sponsored    *SponsoredMark
	Deal         *DealBadge
}

// FeedItem points at a storefront restaurant; the client already has its card data.
type FeedItem struct {
	RestaurantID string
	Reason       string
	Sponsored    *SponsoredMark
	Deal         *DealBadge
}

type Section struct {
	Kind  string
	Title string
	Items []FeedItem
}

type Feed struct {
	GeneratedAt time.Time
	Banners     []Banner
	Sections    []Section
}
