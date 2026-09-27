package service

import (
	"context"
	"reflect"
	"testing"
	"time"

	cuisineentity "github.com/ChristianDenniss/go-data-model/cuisine/entity"
	"github.com/ChristianDenniss/go-data-model/home/entity"
	merchentity "github.com/ChristianDenniss/go-data-model/merchandising/entity"
	money "github.com/ChristianDenniss/go-data-model/money/entity"
	orderentity "github.com/ChristianDenniss/go-data-model/order/entity"
	promoentity "github.com/ChristianDenniss/go-data-model/promotion/entity"
	restaurantentity "github.com/ChristianDenniss/go-data-model/restaurant/entity"
	storefrontentity "github.com/ChristianDenniss/go-data-model/storefront/entity"
)

type fakeCatalog struct{ c storefrontentity.Catalog }

func (f fakeCatalog) LoadCatalog(context.Context, string) (storefrontentity.Catalog, error) {
	return f.c, nil
}

type fakePromos struct{ active []promoentity.ActivePromotion }

func (f fakePromos) ActiveWithTargets(context.Context, time.Time) ([]promoentity.ActivePromotion, error) {
	return f.active, nil
}

type fakeAds struct {
	bySlot map[string][]merchentity.ServedPlacement
}

func (f fakeAds) Serve(_ context.Context, slot, _ string) ([]merchentity.ServedPlacement, error) {
	return f.bySlot[slot], nil
}

func demoCatalog() storefrontentity.Catalog {
	return storefrontentity.Catalog{
		Cuisines: []cuisineentity.Cuisine{{ID: "cu_sushi", Name: "Sushi"}, {ID: "cu_pizza", Name: "Pizza"}},
		Restaurants: []restaurantentity.Restaurant{
			{ID: "koi", CuisineIDs: []string{"cu_sushi"}, Rating: restaurantentity.Rating{Average: 4.7, Count: 1284}},
			{ID: "slice", CuisineIDs: []string{"cu_pizza"}, Rating: restaurantentity.Rating{Average: 4.5, Count: 892}},
			{ID: "stack", Rating: restaurantentity.Rating{Average: 4.4, Count: 2103}},
			{ID: "tiny", Rating: restaurantentity.Rating{Average: 4.9, Count: 3}},
		},
		Orders: []orderentity.Order{{ID: "o1", RestaurantID: "slice"}},
	}
}

func section(feed entity.Feed, kind string) *entity.Section {
	for i := range feed.Sections {
		if feed.Sections[i].Kind == kind {
			return &feed.Sections[i]
		}
	}
	return nil
}

func ids(items []entity.FeedItem) []string {
	out := make([]string, len(items))
	for i, item := range items {
		out[i] = item.RestaurantID
	}
	return out
}

func TestOrganicRailsIgnoreSponsoredPlacements(t *testing.T) {
	catalog := fakeCatalog{demoCatalog()}
	withoutAds := New(catalog, fakePromos{}, fakeAds{})
	withAds := New(catalog, fakePromos{}, fakeAds{bySlot: map[string][]merchentity.ServedPlacement{
		merchentity.SlotHomeRail: {{
			Placement: merchentity.Placement{ID: "pl1", Slot: merchentity.SlotHomeRail, LegacyRestaurantID: "tiny"},
			BidCents:  10_000,
		}},
	}})

	a, err := withoutAds.BuildFeed(context.Background(), "acct")
	if err != nil {
		t.Fatal(err)
	}
	b, err := withAds.BuildFeed(context.Background(), "acct")
	if err != nil {
		t.Fatal(err)
	}

	for _, kind := range []string{entity.SectionPopular, entity.SectionRecommended} {
		if !reflect.DeepEqual(ids(section(a, kind).Items), ids(section(b, kind).Items)) {
			t.Fatalf("%s rail changed when a placement was added", kind)
		}
		for _, item := range section(b, kind).Items {
			if item.Sponsored != nil {
				t.Fatalf("%s rail item %s carries a sponsored mark", kind, item.RestaurantID)
			}
		}
	}

	sponsored := section(b, entity.SectionSponsored)
	if sponsored == nil || len(sponsored.Items) != 1 || sponsored.Items[0].Sponsored == nil {
		t.Fatalf("expected one labelled sponsored item, got %+v", sponsored)
	}
	if section(a, entity.SectionSponsored) != nil {
		t.Fatal("sponsored section should be omitted with no placements")
	}
}

func TestPopularWeightsRatingCount(t *testing.T) {
	got := ids(Popular(demoCatalog(), 10))
	if got[len(got)-1] != "tiny" {
		t.Fatalf("unexpected popular order %v", got)
	}
}

func TestRecommendedFollowsOrderHistory(t *testing.T) {
	items := Recommended(demoCatalog(), 10)
	if items[0].RestaurantID != "slice" || items[0].Reason != "Because you like Pizza" {
		t.Fatalf("expected pizza first with reason, got %+v", items[0])
	}
}

func TestDealsAttachToTargetedRestaurants(t *testing.T) {
	promos := fakePromos{active: []promoentity.ActivePromotion{{
		Promotion: promoentity.Promotion{ID: "promo1", ChannelID: "ch_skip", Name: "Koi pickup", Description: "20% off pickup", Kind: promoentity.KindPercentOff, ValueBPS: 2000, FulfillmentMode: "pickup"},
		Targets:   []promoentity.Target{{ID: "t1", PromotionID: "promo1", LegacyRestaurantID: "koi"}, {ID: "t2", PromotionID: "promo1", LegacyRestaurantID: "ghost"}},
	}}}
	feed, err := New(fakeCatalog{demoCatalog()}, promos, fakeAds{}).BuildFeed(context.Background(), "acct")
	if err != nil {
		t.Fatal(err)
	}

	deals := section(feed, entity.SectionDeals)
	if deals == nil || len(deals.Items) != 1 || deals.Items[0].RestaurantID != "koi" {
		t.Fatalf("expected koi only in deals, got %+v", deals)
	}
	if deals.Items[0].Deal.Label != "20% off pickup" {
		t.Fatalf("unexpected label %q", deals.Items[0].Deal.Label)
	}
	if len(feed.Banners) != 1 || feed.Banners[0].Kind != entity.BannerDeal || feed.Banners[0].Sponsored != nil {
		t.Fatalf("expected one unsponsored deal banner, got %+v", feed.Banners)
	}
}

func TestDealLabel(t *testing.T) {
	cases := map[string]promoentity.Promotion{
		"Free delivery":  {Kind: promoentity.KindFreeDelivery},
		"$5 off":         {Kind: promoentity.KindAmountOff, Value: money.Money{AmountCents: 500}},
		"15% off pickup": {Kind: promoentity.KindPercentOff, ValueBPS: 1500, FulfillmentMode: "pickup"},
		"Tuesday combo":  {Kind: "bundle", Name: "Tuesday combo"},
	}
	for want, p := range cases {
		if got := DealLabel(p); got != want {
			t.Fatalf("want %q, got %q", want, got)
		}
	}
}
