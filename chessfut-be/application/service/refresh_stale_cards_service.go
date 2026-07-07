package service

import (
	"context"
	"time"

	"github.com/fayupable/chessfut-be/application/port/output"
	"github.com/fayupable/chessfut-be/domain"
)

type RefreshStaleCardsService struct {
	chessComClient output.ChessComClientPort
	cardRepository output.CardRepositoryPort
	cache          output.CachePort
}

func NewRefreshStaleCardsService(
	chessComClient output.ChessComClientPort,
	cardRepository output.CardRepositoryPort,
	cache output.CachePort,
) *RefreshStaleCardsService {
	return &RefreshStaleCardsService{chessComClient: chessComClient, cardRepository: cardRepository, cache: cache}
}

func (s *RefreshStaleCardsService) Execute(ctx context.Context, batchSize int) (int, error) {
	stale, err := s.cardRepository.FindStale(ctx, time.Now(), batchSize)
	if err != nil {
		return 0, err
	}

	refreshed := 0
	for _, existing := range stale {
		freshStats, err := s.chessComClient.GetStats(ctx, existing.Player.Username)
		if err != nil {
			continue
		}

		delta := domain.TotalGames(freshStats) - existing.GamesSnapshot
		if delta < minGamesDeltaForRebuild {
			existing.Stats = freshStats
			existing.ExpiresAt = time.Now().Add(cardTTL)
			_ = s.cardRepository.Save(ctx, existing)
			s.refreshCacheIfPresent(ctx, existing)
			refreshed++
			continue
		}

		rebuilt, err := buildCardByType(ctx, s.chessComClient, existing.Player.Username, existing.CardType)
		if err != nil {
			continue
		}
		_ = s.cardRepository.Save(ctx, rebuilt)
		s.refreshCacheIfPresent(ctx, rebuilt)
		refreshed++
	}

	return refreshed, nil
}

// refreshCacheIfPresent keeps Redis in sync with Postgres for cards that are
// already cached (titled/popular players) — without this, a rebuilt card
// would sit correctly in Postgres while GetCard keeps serving the stale
// Redis copy indefinitely, since the read path checks Redis first.
func (s *RefreshStaleCardsService) refreshCacheIfPresent(ctx context.Context, card domain.Card) {
	if _, found, err := s.cache.GetCard(ctx, card.Player.Username); err == nil && found {
		_ = s.cache.SetCard(ctx, card)
	}
}

func buildCardByType(ctx context.Context, client output.ChessComClientPort, username string, cardType domain.CardType) (domain.Card, error) {
	if cardType == domain.CardTypeDetailed {
		return buildDetailedCard(ctx, client, username)
	}
	return buildFastCard(ctx, client, username)
}
