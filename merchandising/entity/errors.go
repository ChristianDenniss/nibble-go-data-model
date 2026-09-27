package entity

import "errors"

var (
	ErrIDRequired       = errors.New("merchandising id is required")
	ErrNotFound         = errors.New("merchandising record not found")
	ErrSlotInvalid      = errors.New("placement.slot is invalid")
	ErrTargetRequired   = errors.New("placement needs a restaurant, place, or brand target")
	ErrStatusInvalid    = errors.New("status is invalid")
	ErrPricingInvalid   = errors.New("campaign.pricing_model is invalid")
	ErrWindowInvalid    = errors.New("campaign.ends_at must be after starts_at")
	ErrEventKindInvalid = errors.New("event.kind must be impression or click")
)
