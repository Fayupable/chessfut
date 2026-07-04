package input

import (
	"context"

	"github.com/fayupable/chessfut-be/domain"
)

type GetFastCardUseCase interface {
	Execute(ctx context.Context, username string) (domain.Card, error)
}

type GetDetailedCardUseCase interface {
	Execute(ctx context.Context, username string) (domain.Card, error)
}
