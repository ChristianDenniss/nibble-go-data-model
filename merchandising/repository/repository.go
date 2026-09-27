package repository

import (
	"context"
	"time"

	"github.com/ChristianDenniss/go-data-model/merchandising/entity"
)

type Repository interface {
	// ActivePlacements returns placements in slot whose campaign is active at `at`.
	// Empty marketID skips market filtering; otherwise campaigns with a null or matching market qualify.
	ActivePlacements(ctx context.Context, slot, marketID string, at time.Time) ([]entity.ServedPlacement, error)
	UpsertAdvertiser(ctx context.Context, a entity.Advertiser) error
	UpsertCampaign(ctx context.Context, c entity.Campaign) error
	UpsertPlacement(ctx context.Context, p entity.Placement) error
	RecordEvent(ctx context.Context, e entity.Event) error
}
