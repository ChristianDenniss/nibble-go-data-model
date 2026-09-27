package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sort"
	"time"

	"github.com/ChristianDenniss/go-data-model/merchandising/entity"
	"github.com/ChristianDenniss/go-data-model/merchandising/repository"
)

// SlotCaps bounds how many placements a slot shows at once.
var SlotCaps = map[string]int{
	entity.SlotHomeBanner:  3,
	entity.SlotHomeRail:    6,
	entity.SlotSearchTop:   2,
	entity.SlotCategoryTop: 2,
}

type Service struct {
	repo repository.Repository
	now  func() time.Time
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

func (s *Service) RecordAdvertiser(ctx context.Context, a entity.Advertiser) error {
	if a.ID == "" {
		return entity.ErrIDRequired
	}
	if a.Status == "" {
		a.Status = entity.AdvertiserActive
	}
	if a.Status != entity.AdvertiserActive && a.Status != entity.AdvertiserSuspended {
		return entity.ErrStatusInvalid
	}
	return s.repo.UpsertAdvertiser(ctx, a)
}

func (s *Service) RecordCampaign(ctx context.Context, c entity.Campaign) error {
	if c.ID == "" || c.AdvertiserID == "" {
		return entity.ErrIDRequired
	}
	if c.Status == "" {
		c.Status = entity.CampaignDraft
	}
	switch c.Status {
	case entity.CampaignDraft, entity.CampaignActive, entity.CampaignPaused, entity.CampaignEnded:
	default:
		return entity.ErrStatusInvalid
	}
	switch c.PricingModel {
	case entity.PricingCPM, entity.PricingCPC, entity.PricingFlat:
	default:
		return entity.ErrPricingInvalid
	}
	if !c.EndsAt.After(c.StartsAt) {
		return entity.ErrWindowInvalid
	}
	return s.repo.UpsertCampaign(ctx, c)
}

func (s *Service) RecordPlacement(ctx context.Context, p entity.Placement) error {
	if p.ID == "" || p.CampaignID == "" {
		return entity.ErrIDRequired
	}
	if _, ok := SlotCaps[p.Slot]; !ok {
		return entity.ErrSlotInvalid
	}
	if p.LegacyRestaurantID == "" && p.PlaceID == "" && p.BrandID == "" {
		return entity.ErrTargetRequired
	}
	return s.repo.UpsertPlacement(ctx, p)
}

// RecordEvent logs an impression or click. ID and OccurredAt are filled when empty.
func (s *Service) RecordEvent(ctx context.Context, e entity.Event) (entity.Event, error) {
	if e.PlacementID == "" {
		return entity.Event{}, entity.ErrIDRequired
	}
	if e.Kind != entity.EventImpression && e.Kind != entity.EventClick {
		return entity.Event{}, entity.ErrEventKindInvalid
	}
	if e.ID == "" {
		e.ID = newEventID()
	}
	if e.OccurredAt.IsZero() {
		e.OccurredAt = s.now().UTC()
	}
	if err := s.repo.RecordEvent(ctx, e); err != nil {
		return entity.Event{}, err
	}
	return e, nil
}

// Serve returns the placements to render in slot right now, ranked and capped.
func (s *Service) Serve(ctx context.Context, slot, marketID string) ([]entity.ServedPlacement, error) {
	if _, ok := SlotCaps[slot]; !ok {
		return nil, entity.ErrSlotInvalid
	}
	live, err := s.repo.ActivePlacements(ctx, slot, marketID, s.now())
	if err != nil {
		return nil, err
	}
	return SelectForSlot(live, slot), nil
}

// SelectForSlot orders by bid, then priority, then id; keeps one placement per target; caps per slot.
func SelectForSlot(live []entity.ServedPlacement, slot string) []entity.ServedPlacement {
	ranked := make([]entity.ServedPlacement, 0, len(live))
	for _, p := range live {
		if p.Placement.Slot == slot {
			ranked = append(ranked, p)
		}
	}
	sort.SliceStable(ranked, func(i, j int) bool {
		a, b := ranked[i], ranked[j]
		if a.BidCents != b.BidCents {
			return a.BidCents > b.BidCents
		}
		if a.Placement.Priority != b.Placement.Priority {
			return a.Placement.Priority > b.Placement.Priority
		}
		return a.Placement.ID < b.Placement.ID
	})

	limit := SlotCaps[slot]
	seen := map[string]bool{}
	out := make([]entity.ServedPlacement, 0, limit)
	for _, p := range ranked {
		key := targetKey(p.Placement)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, p)
		if len(out) == limit {
			break
		}
	}
	return out
}

func targetKey(p entity.Placement) string {
	switch {
	case p.LegacyRestaurantID != "":
		return "restaurant:" + p.LegacyRestaurantID
	case p.PlaceID != "":
		return "place:" + p.PlaceID
	default:
		return "brand:" + p.BrandID
	}
}

func newEventID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return "spe_" + hex.EncodeToString(b[:])
}
