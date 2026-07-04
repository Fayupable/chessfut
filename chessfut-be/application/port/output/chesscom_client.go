package output

import (
	"context"
	"time"

	"github.com/fayupable/chessfut-be/domain"
)

type ChessComClientPort interface {
	GetProfile(ctx context.Context, username string) (domain.Player, error)
	GetStats(ctx context.Context, username string) (domain.PlayerStats, error)
	GetGames(ctx context.Context, username string, from, to time.Time) ([]domain.Game, error)
	GetTitledUsernames(ctx context.Context, title domain.Title) ([]string, error)
	GetLeaderboards(ctx context.Context) (domain.Leaderboards, error)
}
