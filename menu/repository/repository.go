package repository

import (
	"context"

	"github.com/ChristianDenniss/go-data-model/menu/entity"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (entity.Item, error)
	Upsert(ctx context.Context, item entity.Item) error
	// SetImageURL replaces the item's image; empty clears it. Returns entity.ErrNotFound for an unknown id.
	SetImageURL(ctx context.Context, id, imageURL string) error
}
