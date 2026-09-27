package service

import (
	"context"
	"time"

	"github.com/ChristianDenniss/go-data-model/promotion/entity"
	"github.com/ChristianDenniss/go-data-model/promotion/repository"
)

type Service struct {
	promos      repository.PromotionRepository
	memberships repository.MembershipProductRepository
}

func New(promos repository.PromotionRepository, memberships repository.MembershipProductRepository) *Service {
	return &Service{promos: promos, memberships: memberships}
}

func (s *Service) RecordPromotion(ctx context.Context, p entity.Promotion) error {
	if p.ID == "" {
		return entity.ErrIDRequired
	}
	return s.promos.Upsert(ctx, p)
}

func (s *Service) ActivePromotions(ctx context.Context, channelID string, at time.Time) ([]entity.Promotion, error) {
	return s.promos.ListActiveByChannel(ctx, channelID, at)
}

func (s *Service) RecordConstraint(ctx context.Context, c entity.Constraint) error {
	if c.ID == "" || c.PromotionID == "" {
		return entity.ErrIDRequired
	}
	return s.promos.UpsertConstraint(ctx, c)
}

func (s *Service) RecordTarget(ctx context.Context, t entity.Target) error {
	if t.ID == "" || t.PromotionID == "" {
		return entity.ErrIDRequired
	}
	return s.promos.UpsertTarget(ctx, t)
}

func (s *Service) ActiveWithTargets(ctx context.Context, at time.Time) ([]entity.ActivePromotion, error) {
	return s.promos.ListActiveWithTargets(ctx, at)
}

func (s *Service) RecordMembershipProduct(ctx context.Context, p entity.MembershipProduct) error {
	if p.ID == "" {
		return entity.ErrIDRequired
	}
	return s.memberships.Upsert(ctx, p)
}
