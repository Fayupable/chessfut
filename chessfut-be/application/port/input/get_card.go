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

type GetLeaderboardUseCase interface {
	Execute(ctx context.Context, limit, offset int) ([]domain.Card, error)
}

type SearchPlayersUseCase interface {
	Execute(ctx context.Context, query string, limit, offset int) ([]domain.Card, error)
}
