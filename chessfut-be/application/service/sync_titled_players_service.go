package service

import (
	"context"

	"github.com/fayupable/chessfut-be/application/port/input"
	"github.com/fayupable/chessfut-be/application/port/output"
	"github.com/fayupable/chessfut-be/domain"
)

var syncedTitles = []domain.Title{
	domain.TitleGM, domain.TitleWGM,
	domain.TitleIM, domain.TitleWIM,
}

type SyncTitledPlayersService struct {
	chessComClient output.ChessComClientPort
	cardRepository output.CardRepositoryPort
}

func NewSyncTitledPlayersService(
	chessComClient output.ChessComClientPort,
	cardRepository output.CardRepositoryPort,
) *SyncTitledPlayersService {
	return &SyncTitledPlayersService{chessComClient: chessComClient, cardRepository: cardRepository}
}

var _ input.SyncTitledPlayersUseCase = (*SyncTitledPlayersService)(nil)

func (s *SyncTitledPlayersService) Execute(ctx context.Context) (int, error) {
	synced := 0

	for _, title := range syncedTitles {
		usernames, err := s.chessComClient.GetTitledUsernames(ctx, title)
		if err != nil {
			continue
		}

		for _, username := range usernames {
			if _, found, _ := s.cardRepository.FindByUsername(ctx, username); found {
				continue
			}

			card, err := buildFastCard(ctx, s.chessComClient, username)
			if err != nil {
				continue
			}
			if err := s.cardRepository.Save(ctx, card); err != nil {
				continue
			}
			synced++
		}
	}

	return synced, nil
}
