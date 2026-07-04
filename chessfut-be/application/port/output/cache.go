package output

import (
	"context"

	"github.com/fayupable/chessfut-be/domain"
)

type CachePort interface {
	GetCard(ctx context.Context, username string) (domain.Card, bool, error)
	SetCard(ctx context.Context, card domain.Card) error

	IncrementViewCount(ctx context.Context, username string) (int, error)
	IsPromoted(ctx context.Context, username string) (bool, error)
}
