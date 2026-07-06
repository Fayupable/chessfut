package output

import (
	"context"
	"time"

	"github.com/fayupable/chessfut-be/domain"
)

type CardRepositoryPort interface {
	Save(ctx context.Context, card domain.Card) error
	FindByUsername(ctx context.Context, username string) (domain.Card, bool, error)
	FindStale(ctx context.Context, before time.Time, limit int) ([]domain.Card, error)
	FindTopByOVR(ctx context.Context, limit, offset int) ([]domain.Card, error)
	CountAll(ctx context.Context) (int, error)
	SearchByUsername(ctx context.Context, prefix string, limit, offset int) ([]domain.Card, error)
}
