package service

import (
	"context"

	"github.com/fayupable/chessfut-be/application/port/input"
	"github.com/fayupable/chessfut-be/application/port/output"
)

type GetStatsService struct {
	cardRepository output.CardRepositoryPort
	cache          output.CachePort
}

func NewGetStatsService(cardRepository output.CardRepositoryPort, cache output.CachePort) *GetStatsService {
	return &GetStatsService{cardRepository: cardRepository, cache: cache}
}

var _ input.GetStatsUseCase = (*GetStatsService)(nil)

func (s *GetStatsService) Execute(ctx context.Context) (int, error) {
	if count, found, err := s.cache.GetTotalCardsCount(ctx); err == nil && found {
		return count, nil
	}

	count, err := s.cardRepository.CountAll(ctx)
	if err != nil {
		return 0, err
	}

	_ = s.cache.SetTotalCardsCount(ctx, count)
	return count, nil
}
