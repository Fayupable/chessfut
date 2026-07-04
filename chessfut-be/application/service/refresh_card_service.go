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

func (s *RefreshCardService) Execute(ctx context.Context, username string) (domain.Card, error) {
	card, err := buildFastCard(ctx, s.chessComClient, username)
	if err != nil {
		return domain.Card{}, err
	}

	if err := s.cardRepository.Save(ctx, card); err != nil {
		return domain.Card{}, err
	}

	_ = s.cache.SetCard(ctx, card)
	return card, nil
}
