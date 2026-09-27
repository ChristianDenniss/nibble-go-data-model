package repository

import (
	"context"
	"time"

	"github.com/ChristianDenniss/go-data-model/promotion/entity"
)

type PromotionRepository interface {
	GetByID(ctx context.Context, id string) (entity.Promotion, error)
	ListActiveByChannel(ctx context.Context, channelID string, at time.Time) ([]entity.Promotion, error)
	ListActiveWithTargets(ctx context.Context, at time.Time) ([]entity.ActivePromotion, error)
	Upsert(ctx context.Context, p entity.Promotion) error
	UpsertConstraint(ctx context.Context, c entity.Constraint) error
	UpsertTarget(ctx context.Context, t entity.Target) error
}

type MembershipProductRepository interface {
	GetByID(ctx context.Context, id string) (entity.MembershipProduct, error)
	Upsert(ctx context.Context, p entity.MembershipProduct) error
}
