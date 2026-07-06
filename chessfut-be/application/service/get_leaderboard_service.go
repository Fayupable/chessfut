package service

import (
	"context"

	"github.com/fayupable/chessfut-be/application/port/input"
	"github.com/fayupable/chessfut-be/application/port/output"
	"github.com/fayupable/chessfut-be/domain"
)

type GetLeaderboardService struct {
	cardRepository output.CardRepositoryPort
}

func NewGetLeaderboardService(cardRepository output.CardRepositoryPort) *GetLeaderboardService {
	return &GetLeaderboardService{cardRepository: cardRepository}
}

var _ input.GetLeaderboardUseCase = (*GetLeaderboardService)(nil)

func (s *GetLeaderboardService) Execute(ctx context.Context, limit, offset int) ([]domain.Card, error) {
	return s.cardRepository.FindTopByOVR(ctx, limit, offset)
}
