package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/fayupable/chessfut-be/domain"
)

func TestRefreshStaleCardsService_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("touches stats without rebuilding when game delta is below threshold", func(t *testing.T) {
		staleCard := domain.Card{
			Player:        domain.Player{Username: "hikaru"},
			GamesSnapshot: 100,
			CardType:      domain.CardTypeFast,
		}
		freshStats := domain.PlayerStats{Blitz: domain.TimeControlStats{Wins: 60, Losses: 40, Draws: 10}}

		client := new(mockChessComClient)
		repo := new(mockCardRepository)
		cache := new(mockCache)
		repo.On("FindStale", ctx, mock.AnythingOfType("time.Time"), 50).Return([]domain.Card{staleCard}, nil)
		client.On("GetStats", ctx, "hikaru").Return(freshStats, nil)
		repo.On("Save", ctx, mock.AnythingOfType("domain.Card")).Return(nil)

		svc := NewRefreshStaleCardsService(client, repo, cache)
		count, err := svc.Execute(ctx, 50)

		assert.NoError(t, err)
		assert.Equal(t, 1, count)
		client.AssertNotCalled(t, "GetProfile")
		client.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("rebuilds full card when game delta exceeds threshold", func(t *testing.T) {
		staleCard := domain.Card{
			Player:        domain.Player{Username: "hikaru"},
			GamesSnapshot: 0,
			CardType:      domain.CardTypeFast,
		}
		freshStats := domain.PlayerStats{Blitz: domain.TimeControlStats{Wins: 60, Losses: 40, Draws: 10}}
		player := domain.Player{Username: "hikaru", Title: domain.TitleGM}

		client := new(mockChessComClient)
		repo := new(mockCardRepository)
		cache := new(mockCache)
		repo.On("FindStale", ctx, mock.AnythingOfType("time.Time"), 50).Return([]domain.Card{staleCard}, nil)
		client.On("GetStats", ctx, "hikaru").Return(freshStats, nil).Once()
		client.On("GetProfile", ctx, "hikaru").Return(player, nil)
		client.On("GetStats", ctx, "hikaru").Return(freshStats, nil)
		repo.On("Save", ctx, mock.AnythingOfType("domain.Card")).Return(nil)

		svc := NewRefreshStaleCardsService(client, repo, cache)
		count, err := svc.Execute(ctx, 50)

		assert.NoError(t, err)
		assert.Equal(t, 1, count)
		repo.AssertExpectations(t)
	})
}

var _ = time.Now

func TestRefreshStaleCardsService_Execute_RebuildsDetailedCardType(t *testing.T) {
	ctx := context.Background()
	staleCard := domain.Card{
		Player:        domain.Player{Username: "hikaru"},
		GamesSnapshot: 0,
		CardType:      domain.CardTypeDetailed,
	}
	freshStats := domain.PlayerStats{Blitz: domain.TimeControlStats{Wins: 60, Losses: 40, Draws: 10}}
	player := domain.Player{Username: "hikaru", Title: domain.TitleGM}
	games := []domain.Game{{MovesCount: 30, Result: domain.ResultWin}}

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	repo.On("FindStale", ctx, mock.AnythingOfType("time.Time"), 50).Return([]domain.Card{staleCard}, nil)
	client.On("GetStats", ctx, "hikaru").Return(freshStats, nil).Once()
	client.On("GetProfile", ctx, "hikaru").Return(player, nil)
	client.On("GetStats", ctx, "hikaru").Return(freshStats, nil)
	client.On("GetGames", ctx, "hikaru", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(games, nil)
	repo.On("Save", ctx, mock.AnythingOfType("domain.Card")).Return(nil)

	svc := NewRefreshStaleCardsService(client, repo, cache)
	count, err := svc.Execute(ctx, 50)

	assert.NoError(t, err)
	assert.Equal(t, 1, count)
	client.AssertExpectations(t)
	repo.AssertExpectations(t)
}
func TestRefreshStaleCardsService_Execute_ReturnsErrorWhenFindStaleFails(t *testing.T) {
	ctx := context.Background()

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	repo.On("FindStale", ctx, mock.AnythingOfType("time.Time"), 50).Return([]domain.Card(nil), assert.AnError)

	svc := NewRefreshStaleCardsService(client, repo, cache)
	_, err := svc.Execute(ctx, 50)

	assert.Error(t, err)
}

func TestRefreshStaleCardsService_Execute_SkipsCardWhenStatsFetchFails(t *testing.T) {
	ctx := context.Background()
	staleCard := domain.Card{Player: domain.Player{Username: "hikaru"}, CardType: domain.CardTypeFast}

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	repo.On("FindStale", ctx, mock.AnythingOfType("time.Time"), 50).Return([]domain.Card{staleCard}, nil)
	client.On("GetStats", ctx, "hikaru").Return(domain.PlayerStats{}, assert.AnError)

	svc := NewRefreshStaleCardsService(client, repo, cache)
	count, err := svc.Execute(ctx, 50)

	assert.NoError(t, err)
	assert.Equal(t, 0, count)
	repo.AssertNotCalled(t, "Save")
}
