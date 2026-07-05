package service

import (
	"context"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/fayupable/chessfut-be/application/port/input"
	"github.com/fayupable/chessfut-be/application/port/output"
	"github.com/fayupable/chessfut-be/domain"
)

const detailedGamesWindow = 30 * 24 * time.Hour
const minGamesDeltaForRebuild = 50

type GetDetailedCardService struct {
	chessComClient output.ChessComClientPort
	cardRepository output.CardRepositoryPort
	cache          output.CachePort
	group          singleflight.Group
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
	if card, ok := s.tryCache(ctx, username); ok {
		return card, nil
	}

	existing, found, err := s.cardRepository.FindByUsername(ctx, username)
	if err != nil {
		return domain.Card{}, err
	}
	if found && existing.CardType == domain.CardTypeDetailed {
		s.promoteIfPopular(ctx, username, existing)
		return existing, nil
	}

	return s.rebuild(ctx, username)
}

func (s *GetDetailedCardService) tryCache(ctx context.Context, username string) (domain.Card, bool) {
	card, found, err := s.cache.GetCard(ctx, username)
	if err != nil || !found || card.CardType != domain.CardTypeDetailed {
		return domain.Card{}, false
	}
	return card, true
}

// rebuild is guarded by singleflight so concurrent requests for the same
// never-before-seen username (e.g. from the homepage card fan and a direct
// page visit firing at once) share a single chess.com fetch instead of each
// independently hammering the API.
func (s *GetDetailedCardService) rebuild(ctx context.Context, username string) (domain.Card, error) {
	result, err, _ := s.group.Do(username, func() (any, error) {
		card, err := buildDetailedCard(ctx, s.chessComClient, username)
		if err != nil {
			return domain.Card{}, err
		}
		if err := s.cardRepository.Save(ctx, card); err != nil {
			return domain.Card{}, err
		}
		return card, nil
	})
	if err != nil {
		return domain.Card{}, err
	}

	card := result.(domain.Card)
	s.promoteIfPopular(ctx, username, card)
	return card, nil
}

func (s *GetDetailedCardService) promoteIfPopular(ctx context.Context, username string, card domain.Card) {
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

func buildDetailedCard(ctx context.Context, client output.ChessComClientPort, username string) (domain.Card, error) {
	player, err := client.GetProfile(ctx, username)
	if err != nil {
		return domain.Card{}, err
	}

	stats, err := client.GetStats(ctx, username)
	if err != nil {
		return domain.Card{}, err
	}

	to := time.Now()
	from := to.Add(-detailedGamesWindow)
	games, err := client.GetGames(ctx, username, from, to)
	if err != nil {
		return domain.Card{}, err
	}

	tier := domain.CardTierStandard
	if player.HasFideTitle() {
		tier = domain.CardTierTitled
	}

	playStyle := DeterminePlayStyle(games)
	topOpenings := aggregateTopOpenings(games)
	gamesSnapshot := domain.TotalGames(stats)
	attributes, position, ovr, workRate := BuildCardScoring(player, stats, topOpenings, games, gamesSnapshot)

	now := time.Now()
	return domain.Card{
		Player:        player,
		Stats:         stats,
		CardType:      domain.CardTypeDetailed,
		Tier:          tier,
		OVR:           ovr,
		PlayStyle:     playStyle,
		Position:      position,
		Attributes:    attributes,
		WorkRate:      workRate,
		Badges:        AssignBadges(player, stats),
		TopOpenings:   topOpenings,
		GamesSnapshot: gamesSnapshot,
		ComputedAt:    now,
		ExpiresAt:     now.Add(cardTTL),
	}, nil
}
