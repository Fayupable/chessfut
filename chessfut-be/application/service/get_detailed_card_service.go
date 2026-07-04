package service

import (
	"context"
	"time"

	"github.com/fayupable/chessfut-be/application/port/input"
	"github.com/fayupable/chessfut-be/application/port/output"
	"github.com/fayupable/chessfut-be/domain"
)

const detailedGamesWindow = 2 * 30 * 24 * time.Hour

type GetDetailedCardService struct {
	chessComClient output.ChessComClientPort
	cardRepository output.CardRepositoryPort
	cache          output.CachePort
}

func NewGetDetailedCardService(
	chessComClient output.ChessComClientPort,
	cardRepository output.CardRepositoryPort,
	cache output.CachePort,
) *GetDetailedCardService {
	return &GetDetailedCardService{
		chessComClient: chessComClient,
		cardRepository: cardRepository,
		cache:          cache,
	}
}

var _ input.GetDetailedCardUseCase = (*GetDetailedCardService)(nil)

func (s *GetDetailedCardService) Execute(ctx context.Context, username string) (domain.Card, error) {
	if card, found, err := s.cache.GetCard(ctx, username); err == nil && found && !card.IsExpired() && card.CardType == domain.CardTypeDetailed {
		return card, nil
	}

	if card, found, err := s.cardRepository.FindByUsername(ctx, username); err == nil && found && !card.IsExpired() && card.CardType == domain.CardTypeDetailed {
		return card, nil
	}

	card, err := s.buildDetailedCard(ctx, username)
	if err != nil {
		return domain.Card{}, err
	}

	if err := s.cardRepository.Save(ctx, card); err != nil {
		return domain.Card{}, err
	}

	return card, nil
}

func (s *GetDetailedCardService) buildDetailedCard(ctx context.Context, username string) (domain.Card, error) {
	player, err := s.chessComClient.GetProfile(ctx, username)
	if err != nil {
		return domain.Card{}, err
	}

	stats, err := s.chessComClient.GetStats(ctx, username)
	if err != nil {
		return domain.Card{}, err
	}

	to := time.Now()
	from := to.Add(-detailedGamesWindow)
	games, err := s.chessComClient.GetGames(ctx, username, from, to)
	if err != nil {
		return domain.Card{}, err
	}

	tier := domain.CardTierStandard
	if player.HasFideTitle() {
		tier = domain.CardTierTitled
	}

	now := time.Now()
	return domain.Card{
		Player:      player,
		Stats:       stats,
		CardType:    domain.CardTypeDetailed,
		Tier:        tier,
		OVR:         CalculateOVR(stats),
		PlayStyle:   DeterminePlayStyle(games),
		Badges:      AssignBadges(player, stats),
		TopOpenings: aggregateTopOpenings(games),
		ComputedAt:  now,
		ExpiresAt:   now.Add(cardTTL),
	}, nil
}
