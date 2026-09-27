package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ChristianDenniss/go-data-model/home/entity"
	merchentity "github.com/ChristianDenniss/go-data-model/merchandising/entity"
	promoentity "github.com/ChristianDenniss/go-data-model/promotion/entity"
	storefrontentity "github.com/ChristianDenniss/go-data-model/storefront/entity"
)

const (
	railLimit   = 10
	// Keep the sponsored placements and five organic deal banners visible.
	bannerLimit = 8
)

type CatalogLoader interface {
	LoadCatalog(ctx context.Context, accountID string) (storefrontentity.Catalog, error)
}

type PromotionSource interface {
	ActiveWithTargets(ctx context.Context, at time.Time) ([]promoentity.ActivePromotion, error)
}

type PlacementSource interface {
	Serve(ctx context.Context, slot, marketID string) ([]merchentity.ServedPlacement, error)
}

type Service struct {
	catalog CatalogLoader
	promos  PromotionSource
	ads     PlacementSource
	now     func() time.Time
}

func New(catalog CatalogLoader, promos PromotionSource, ads PlacementSource) *Service {
	return &Service{catalog: catalog, promos: promos, ads: ads, now: time.Now}
}

// BuildFeed composes Home. Organic rails are computed from the catalog alone, before
// merchandising is read, so paid placement cannot reorder them (D28).
func (s *Service) BuildFeed(ctx context.Context, accountID string) (entity.Feed, error) {
	catalog, err := s.catalog.LoadCatalog(ctx, accountID)
	if err != nil {
		return entity.Feed{}, err
	}
	at := s.now().UTC()

	popular := Popular(catalog, railLimit)
	recommended := Recommended(catalog, railLimit)

	active, err := s.promos.ActiveWithTargets(ctx, at)
	if err != nil {
		return entity.Feed{}, err
	}
	known := restaurantSet(catalog)
	deals := dealsByRestaurant(active, known)

	banners, err := s.ads.Serve(ctx, merchentity.SlotHomeBanner, "")
	if err != nil {
		return entity.Feed{}, err
	}
	rail, err := s.ads.Serve(ctx, merchentity.SlotHomeRail, "")
	if err != nil {
		return entity.Feed{}, err
	}

	feed := entity.Feed{GeneratedAt: at}
	feed.Banners = buildBanners(banners, active, known, deals)

	if items := sponsoredItems(rail, known, deals); len(items) > 0 {
		feed.Sections = append(feed.Sections, entity.Section{Kind: entity.SectionSponsored, Title: "Featured near you", Items: items})
	}
	if items := dealItems(catalog, deals); len(items) > 0 {
		feed.Sections = append(feed.Sections, entity.Section{Kind: entity.SectionDeals, Title: "Popular deals in your area", Items: items})
	}
	feed.Sections = append(feed.Sections,
		entity.Section{Kind: entity.SectionPopular, Title: "Most popular", Items: withDeals(popular, deals)},
		entity.Section{Kind: entity.SectionRecommended, Title: "Recommended for you", Items: withDeals(recommended, deals)},
	)
	return feed, nil
}

func restaurantSet(catalog storefrontentity.Catalog) map[string]bool {
	out := make(map[string]bool, len(catalog.Restaurants))
	for _, r := range catalog.Restaurants {
		out[r.ID] = true
	}
	return out
}

func dealsByRestaurant(active []promoentity.ActivePromotion, known map[string]bool) map[string]*entity.DealBadge {
	out := map[string]*entity.DealBadge{}
	for _, ap := range active {
		for _, t := range ap.Targets {
			if t.LegacyRestaurantID == "" || !known[t.LegacyRestaurantID] {
				continue
			}
			if _, taken := out[t.LegacyRestaurantID]; taken {
				continue
			}
			out[t.LegacyRestaurantID] = badgeFor(ap.Promotion)
		}
	}
	return out
}

func badgeFor(p promoentity.Promotion) *entity.DealBadge {
	return &entity.DealBadge{
		PromotionID:     p.ID,
		ChannelID:       p.ChannelID,
		Label:           DealLabel(p),
		FulfillmentMode: p.FulfillmentMode,
	}
}

// DealLabel is the short badge copy for a promotion, e.g. "20% off pickup".
func DealLabel(p promoentity.Promotion) string {
	var label string
	switch p.Kind {
	case promoentity.KindFreeDelivery:
		return "Free delivery"
	case promoentity.KindPercentOff:
		label = fmt.Sprintf("%d%% off", p.ValueBPS/100)
	case promoentity.KindAmountOff:
		label = fmt.Sprintf("$%d off", p.Value.AmountCents/100)
	default:
		if p.Name != "" {
			return p.Name
		}
		return "Deal"
	}
	switch p.FulfillmentMode {
	case "pickup":
		label += " pickup"
	case "delivery":
		label += " delivery"
	}
	return label
}

func sponsoredMark(p merchentity.ServedPlacement) *entity.SponsoredMark {
	return &entity.SponsoredMark{
		PlacementID:    p.Placement.ID,
		CampaignID:     p.CampaignID,
		AdvertiserName: p.AdvertiserName,
		Label:          "Sponsored",
	}
}

func buildBanners(
	served []merchentity.ServedPlacement,
	active []promoentity.ActivePromotion,
	known map[string]bool,
	deals map[string]*entity.DealBadge,
) []entity.Banner {
	out := make([]entity.Banner, 0, bannerLimit)
	for _, p := range served {
		restaurantID := p.Placement.LegacyRestaurantID
		if restaurantID != "" && !known[restaurantID] {
			continue
		}
		out = append(out, entity.Banner{
			ID:           p.Placement.ID,
			Kind:         entity.BannerSponsored,
			Headline:     p.Placement.Headline,
			Body:         p.Placement.Body,
			ImageURL:     p.Placement.ImageURL,
			CallToAction: p.Placement.CallToAction,
			RestaurantID: restaurantID,
			Sponsored:    sponsoredMark(p),
			Deal:         deals[restaurantID],
		})
	}
	for _, ap := range active {
		if len(out) == bannerLimit {
			break
		}
		if ap.Promotion.Description == "" {
			continue
		}
		restaurantID := ""
		for _, t := range ap.Targets {
			if known[t.LegacyRestaurantID] {
				restaurantID = t.LegacyRestaurantID
				break
			}
		}
		out = append(out, entity.Banner{
			ID:           ap.Promotion.ID,
			Kind:         entity.BannerDeal,
			Headline:     ap.Promotion.Name,
			Body:         ap.Promotion.Description,
			CallToAction: "See deal",
			RestaurantID: restaurantID,
			Deal:         badgeFor(ap.Promotion),
		})
	}
	return out
}

func sponsoredItems(served []merchentity.ServedPlacement, known map[string]bool, deals map[string]*entity.DealBadge) []entity.FeedItem {
	out := make([]entity.FeedItem, 0, len(served))
	for _, p := range served {
		id := p.Placement.LegacyRestaurantID
		if id == "" || !known[id] {
			continue
		}
		out = append(out, entity.FeedItem{RestaurantID: id, Reason: p.Placement.Headline, Sponsored: sponsoredMark(p), Deal: deals[id]})
	}
	return out
}

// dealItems lists restaurants with an active deal, most popular first.
func dealItems(catalog storefrontentity.Catalog, deals map[string]*entity.DealBadge) []entity.FeedItem {
	out := make([]entity.FeedItem, 0, len(deals))
	for _, item := range Popular(catalog, len(catalog.Restaurants)) {
		if d := deals[item.RestaurantID]; d != nil {
			out = append(out, entity.FeedItem{RestaurantID: item.RestaurantID, Deal: d})
		}
	}
	return out
}

func withDeals(items []entity.FeedItem, deals map[string]*entity.DealBadge) []entity.FeedItem {
	out := make([]entity.FeedItem, len(items))
	for i, item := range items {
		item.Deal = deals[item.RestaurantID]
		out[i] = item
	}
	return out
}
