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
			refreshed++
			continue
		}

		rebuilt, err := buildCardByType(ctx, s.chessComClient, existing.Player.Username, existing.CardType)
		if err != nil {
			continue
		}
		_ = s.cardRepository.Save(ctx, rebuilt)
		refreshed++
	}

	return refreshed, nil
}

func buildCardByType(ctx context.Context, client output.ChessComClientPort, username string, cardType domain.CardType) (domain.Card, error) {
	if cardType == domain.CardTypeDetailed {
		return buildDetailedCard(ctx, client, username)
	}
	return buildFastCard(ctx, client, username)
}
