package service

import (
	"context"

	"github.com/fayupable/chessfut-be/application/port/input"
	"github.com/fayupable/chessfut-be/application/port/output"
	"github.com/fayupable/chessfut-be/domain"
)

type SearchPlayersService struct {
	cardRepository output.CardRepositoryPort
}

func NewSearchPlayersService(cardRepository output.CardRepositoryPort) *SearchPlayersService {
	return &SearchPlayersService{cardRepository: cardRepository}
}

var _ input.SearchPlayersUseCase = (*SearchPlayersService)(nil)

func (s *SearchPlayersService) Execute(ctx context.Context, query string, limit, offset int) ([]domain.Card, error) {
	return s.cardRepository.SearchByUsername(ctx, query, limit, offset)
}
