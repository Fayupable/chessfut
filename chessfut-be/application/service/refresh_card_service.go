package service

import (
	"context"

	"github.com/fayupable/chessfut-be/application/port/input"
	"github.com/fayupable/chessfut-be/application/port/output"
	"github.com/fayupable/chessfut-be/domain"
)

type RefreshCardService struct {
	chessComClient output.ChessComClientPort
	cardRepository output.CardRepositoryPort
	cache          output.CachePort
}

func NewRefreshCardService(
	chessComClient output.ChessComClientPort,
	cardRepository output.CardRepositoryPort,
	cache output.CachePort,
) *RefreshCardService {
	return &RefreshCardService{chessComClient: chessComClient, cardRepository: cardRepository, cache: cache}
}

var _ input.RefreshCardUseCase = (*RefreshCardService)(nil)

// Execute rebuilds the card as whatever type it already was (fast or
// detailed) — defaulting to fast only when no card exists yet — so an admin
// refresh can never silently downgrade a detailed card to fast.
func (s *RefreshCardService) Execute(ctx context.Context, username string) (domain.Card, error) {
	cardType := domain.CardTypeFast
	if existing, found, err := s.cardRepository.FindByUsername(ctx, username); err == nil && found {
		cardType = existing.CardType
	}

	card, err := buildCardByType(ctx, s.chessComClient, username, cardType)
	if err != nil {
		return domain.Card{}, err
	}

	if err := s.cardRepository.Save(ctx, card); err != nil {
		return domain.Card{}, err
	}

	_ = s.cache.SetCard(ctx, card)
	return card, nil
}
