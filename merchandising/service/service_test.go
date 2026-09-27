package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ChristianDenniss/go-data-model/merchandising/entity"
)

func served(id, restaurant string, bid int64, priority int) entity.ServedPlacement {
	return entity.ServedPlacement{
		Placement: entity.Placement{ID: id, Slot: entity.SlotHomeRail, LegacyRestaurantID: restaurant, Priority: priority},
		BidCents:  bid,
	}
}

func TestSelectForSlotRanksByBidThenPriority(t *testing.T) {
	got := SelectForSlot([]entity.ServedPlacement{
		served("a", "r1", 100, 0),
		served("b", "r2", 300, 0),
		served("c", "r3", 100, 5),
	}, entity.SlotHomeRail)

	want := []string{"b", "c", "a"}
	for i, id := range want {
		if got[i].Placement.ID != id {
			t.Fatalf("position %d: want %s, got %s", i, id, got[i].Placement.ID)
		}
	}
}

func TestSelectForSlotDedupesTargetAndCaps(t *testing.T) {
	live := []entity.ServedPlacement{served("dup-low", "r1", 50, 0), served("dup-high", "r1", 500, 0)}
	for i := 0; i < 10; i++ {
		live = append(live, served(string(rune('k'+i)), string(rune('A'+i)), 10, 0))
	}
	got := SelectForSlot(live, entity.SlotHomeRail)

	if len(got) != SlotCaps[entity.SlotHomeRail] {
		t.Fatalf("want cap %d, got %d", SlotCaps[entity.SlotHomeRail], len(got))
	}
	if got[0].Placement.ID != "dup-high" {
		t.Fatalf("highest bid for r1 should win, got %s", got[0].Placement.ID)
	}
	for _, p := range got[1:] {
		if p.Placement.LegacyRestaurantID == "r1" {
			t.Fatal("r1 served twice")
		}
	}
}

func TestSelectForSlotIgnoresOtherSlots(t *testing.T) {
	banner := served("banner", "r1", 999, 0)
	banner.Placement.Slot = entity.SlotHomeBanner
	got := SelectForSlot([]entity.ServedPlacement{banner, served("rail", "r2", 1, 0)}, entity.SlotHomeRail)
	if len(got) != 1 || got[0].Placement.ID != "rail" {
		t.Fatalf("expected only the rail placement, got %+v", got)
	}
}

type nopRepo struct{ events []entity.Event }

func (nopRepo) ActivePlacements(context.Context, string, string, time.Time) ([]entity.ServedPlacement, error) {
	return nil, nil
}
func (nopRepo) UpsertAdvertiser(context.Context, entity.Advertiser) error { return nil }
func (nopRepo) UpsertCampaign(context.Context, entity.Campaign) error     { return nil }
func (nopRepo) UpsertPlacement(context.Context, entity.Placement) error   { return nil }
func (r *nopRepo) RecordEvent(_ context.Context, e entity.Event) error {
	r.events = append(r.events, e)
	return nil
}

func TestRecordCampaignValidatesWindowAndPricing(t *testing.T) {
	svc := New(&nopRepo{})
	start := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)

	err := svc.RecordCampaign(context.Background(), entity.Campaign{ID: "c", AdvertiserID: "a", PricingModel: entity.PricingCPM, StartsAt: start, EndsAt: start})
	if !errors.Is(err, entity.ErrWindowInvalid) {
		t.Fatalf("want ErrWindowInvalid, got %v", err)
	}
	err = svc.RecordCampaign(context.Background(), entity.Campaign{ID: "c", AdvertiserID: "a", PricingModel: "cpa", StartsAt: start, EndsAt: start.Add(time.Hour)})
	if !errors.Is(err, entity.ErrPricingInvalid) {
		t.Fatalf("want ErrPricingInvalid, got %v", err)
	}
}

func TestRecordPlacementRequiresTarget(t *testing.T) {
	svc := New(&nopRepo{})
	err := svc.RecordPlacement(context.Background(), entity.Placement{ID: "p", CampaignID: "c", Slot: entity.SlotHomeRail})
	if !errors.Is(err, entity.ErrTargetRequired) {
		t.Fatalf("want ErrTargetRequired, got %v", err)
	}
}

func TestRecordEventFillsIDAndTime(t *testing.T) {
	repo := &nopRepo{}
	svc := New(repo)
	got, err := svc.RecordEvent(context.Background(), entity.Event{PlacementID: "p", Kind: entity.EventClick})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID == "" || got.OccurredAt.IsZero() || len(repo.events) != 1 {
		t.Fatalf("event not filled/recorded: %+v", got)
	}
	if _, err := svc.RecordEvent(context.Background(), entity.Event{PlacementID: "p", Kind: "view"}); !errors.Is(err, entity.ErrEventKindInvalid) {
		t.Fatalf("want ErrEventKindInvalid, got %v", err)
	}
}
