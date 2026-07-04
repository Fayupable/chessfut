package input

import (
	"context"

	"github.com/fayupable/chessfut-be/domain"
)

type RefreshCardUseCase interface {
	Execute(ctx context.Context, username string) (domain.Card, error)
}

type SyncTitledPlayersUseCase interface {
	Execute(ctx context.Context) (int, error)
}
