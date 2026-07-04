package service

import (
	"context"
	"time"

	"github.com/fayupable/chessfut-be/application/port/input"
	"github.com/fayupable/chessfut-be/application/port/output"
	"github.com/fayupable/chessfut-be/domain"
)

const cardTTL = 48 * time.Hour
const viewPromotionThreshold = 5

type GetFastCardService struct {
	chessComClient output.ChessComClientPort
	cardRepository output.CardRepositoryPort
	cache          output.CachePort
}

func NewGetFastCardService(
	chessComClient output.ChessComClientPort,
	cardRepository output.CardRepositoryPort,
	cache output.CachePort,
) *GetFastCardService {
	return &GetFastCardService{
		chessComClient: chessComClient,
		cardRepository: cardRepository,
		cache:          cache,
	}
}

var _ input.GetFastCardUseCase = (*GetFastCardService)(nil)

func (s *GetFastCardService) Execute(ctx context.Context, username string) (domain.Card, error) {
	if card, found, err := s.cache.GetCard(ctx, username); err == nil && found {
		return card, nil
	}

	if card, found, err := s.cardRepository.FindByUsername(ctx, username); err == nil && found {
		s.promoteIfPopular(ctx, username, card)
		return card, nil
	}

	card, err := buildFastCard(ctx, s.chessComClient, username)
	if err != nil {
		return domain.Card{}, err
	}

	if err := s.cardRepository.Save(ctx, card); err != nil {
		return domain.Card{}, err
	}

	s.promoteIfPopular(ctx, username, card)
	return card, nil
}

func buildFastCard(ctx context.Context, client output.ChessComClientPort, username string) (domain.Card, error) {
	player, err := client.GetProfile(ctx, username)
	if err != nil {
		return domain.Card{}, err
	}

	stats, err := client.GetStats(ctx, username)
	if err != nil {
		return domain.Card{}, err
	}

	tier := domain.CardTierStandard
	if player.HasFideTitle() {
		tier = domain.CardTierTitled
	}

	now := time.Now()
	return domain.Card{
		Player:        player,
		Stats:         stats,
		CardType:      domain.CardTypeFast,
		Tier:          tier,
		OVR:           CalculateOVR(stats),
		Position:      domain.PositionAllRounder,
		Badges:        AssignBadges(player, stats),
		GamesSnapshot: domain.TotalGames(stats),
		ComputedAt:    now,
		ExpiresAt:     now.Add(cardTTL),
	}, nil
}

func (s *GetFastCardService) promoteIfPopular(ctx context.Context, username string, card domain.Card) {
	if card.Player.HasFideTitle() {
		_ = s.cache.SetCard(ctx, card)
		return
	}

	count, err := s.cache.IncrementViewCount(ctx, username)
	if err != nil {
		return
	}
	if count >= viewPromotionThreshold {
		_ = s.cache.SetCard(ctx, card)
	}
}
