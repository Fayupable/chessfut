package output

import (
	"context"

	"github.com/fayupable/chessfut-be/domain"
)

type CardRepositoryPort interface {
	Save(ctx context.Context, card domain.Card) error
	FindByUsername(ctx context.Context, username string) (domain.Card, bool, error)
}
