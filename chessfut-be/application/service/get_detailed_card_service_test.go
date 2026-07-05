package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/fayupable/chessfut-be/domain"
)

func TestGetDetailedCardService_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("returns detailed card from cache without other calls", func(t *testing.T) {
		cachedCard := domain.Card{Player: domain.Player{Username: "hikaru"}, CardType: domain.CardTypeDetailed}

		client := new(mockChessComClient)
		repo := new(mockCardRepository)
		cache := new(mockCache)
		cache.On("GetCard", ctx, "hikaru").Return(cachedCard, true, nil)

		svc := NewGetDetailedCardService(client, repo, cache)
		result, err := svc.Execute(ctx, "hikaru")

		assert.NoError(t, err)
		assert.Equal(t, cachedCard, result)
		repo.AssertNotCalled(t, "FindByUsername")
		client.AssertNotCalled(t, "GetProfile")
	})

	t.Run("returns detailed card from repository and promotes it", func(t *testing.T) {
		existingCard := domain.Card{
			Player:   domain.Player{Username: "hikaru", Title: domain.TitleGM},
			CardType: domain.CardTypeDetailed,
		}

		client := new(mockChessComClient)
		repo := new(mockCardRepository)
		cache := new(mockCache)
		cache.On("GetCard", ctx, "hikaru").Return(domain.Card{}, false, nil)
		repo.On("FindByUsername", ctx, "hikaru").Return(existingCard, true, nil)
		cache.On("SetCard", ctx, existingCard).Return(nil)

		svc := NewGetDetailedCardService(client, repo, cache)
		result, err := svc.Execute(ctx, "hikaru")

		assert.NoError(t, err)
		assert.Equal(t, existingCard, result)
		client.AssertNotCalled(t, "GetProfile")
	})

	t.Run("rebuilds when repository has no detailed card yet", func(t *testing.T) {
		player := domain.Player{Username: "hikaru", Title: domain.TitleGM}
		stats := domain.PlayerStats{FideRating: 2814, Blitz: domain.TimeControlStats{Rating: 3414}, Rapid: domain.TimeControlStats{Rating: 2839}}
		games := []domain.Game{{MovesCount: 30, Result: domain.ResultWin}}

		client := new(mockChessComClient)
		repo := new(mockCardRepository)
		cache := new(mockCache)
		cache.On("GetCard", ctx, "hikaru").Return(domain.Card{}, false, nil)
		repo.On("FindByUsername", ctx, "hikaru").Return(domain.Card{}, false, nil)
		client.On("GetProfile", ctx, "hikaru").Return(player, nil)
		client.On("GetStats", ctx, "hikaru").Return(stats, nil)
		client.On("GetGames", ctx, "hikaru", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(games, nil)
		repo.On("Save", ctx, mock.AnythingOfType("domain.Card")).Return(nil)
		cache.On("SetCard", ctx, mock.AnythingOfType("domain.Card")).Return(nil)

		svc := NewGetDetailedCardService(client, repo, cache)
		result, err := svc.Execute(ctx, "hikaru")

		assert.NoError(t, err)
		assert.Equal(t, domain.CardTypeDetailed, result.CardType)
		assert.Greater(t, result.OVR, 0)
		assert.LessOrEqual(t, result.OVR, 99)
		client.AssertExpectations(t)
		repo.AssertExpectations(t)
	})
}

var _ = time.Now

func TestGetDetailedCardService_Execute_PromotesNonTitledPlayerAtThreshold(t *testing.T) {
	ctx := context.Background()
	player := domain.Player{Username: "someuser", Title: domain.TitleNone}
	stats := domain.PlayerStats{Blitz: domain.TimeControlStats{Rating: 1500}}
	games := []domain.Game{{MovesCount: 30, Result: domain.ResultWin}}

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	cache.On("GetCard", ctx, "someuser").Return(domain.Card{}, false, nil)
	repo.On("FindByUsername", ctx, "someuser").Return(domain.Card{}, false, nil)
	client.On("GetProfile", ctx, "someuser").Return(player, nil)
	client.On("GetStats", ctx, "someuser").Return(stats, nil)
	client.On("GetGames", ctx, "someuser", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(games, nil)
	repo.On("Save", ctx, mock.AnythingOfType("domain.Card")).Return(nil)
	cache.On("IncrementViewCount", ctx, "someuser").Return(5, nil)
	cache.On("SetCard", ctx, mock.AnythingOfType("domain.Card")).Return(nil)

	svc := NewGetDetailedCardService(client, repo, cache)
	_, err := svc.Execute(ctx, "someuser")

	assert.NoError(t, err)
	cache.AssertExpectations(t)
}

func TestGetDetailedCardService_Execute_ReturnsErrorWhenGamesFetchFails(t *testing.T) {
	ctx := context.Background()
	player := domain.Player{Username: "hikaru", Title: domain.TitleGM}
	stats := domain.PlayerStats{FideRating: 2814}

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	cache.On("GetCard", ctx, "hikaru").Return(domain.Card{}, false, nil)
	repo.On("FindByUsername", ctx, "hikaru").Return(domain.Card{}, false, nil)
	client.On("GetProfile", ctx, "hikaru").Return(player, nil)
	client.On("GetStats", ctx, "hikaru").Return(stats, nil)
	client.On("GetGames", ctx, "hikaru", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return([]domain.Game(nil), assert.AnError)

	svc := NewGetDetailedCardService(client, repo, cache)
	_, err := svc.Execute(ctx, "hikaru")

	assert.Error(t, err)
	repo.AssertNotCalled(t, "Save")
}

func TestGetDetailedCardService_Execute_ReturnsErrorWhenRepositoryFails(t *testing.T) {
	ctx := context.Background()

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	cache.On("GetCard", ctx, "hikaru").Return(domain.Card{}, false, nil)
	repo.On("FindByUsername", ctx, "hikaru").Return(domain.Card{}, false, assert.AnError)

	svc := NewGetDetailedCardService(client, repo, cache)
	_, err := svc.Execute(ctx, "hikaru")

	assert.Error(t, err)
	client.AssertNotCalled(t, "GetProfile")
}

func TestGetDetailedCardService_Rebuild_ReturnsErrorWhenSaveFails(t *testing.T) {
	ctx := context.Background()
	player := domain.Player{Username: "hikaru", Title: domain.TitleGM}
	stats := domain.PlayerStats{FideRating: 2814}
	games := []domain.Game{{MovesCount: 30, Result: domain.ResultWin}}

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	cache.On("GetCard", ctx, "hikaru").Return(domain.Card{}, false, nil)
	repo.On("FindByUsername", ctx, "hikaru").Return(domain.Card{}, false, nil)
	client.On("GetProfile", ctx, "hikaru").Return(player, nil)
	client.On("GetStats", ctx, "hikaru").Return(stats, nil)
	client.On("GetGames", ctx, "hikaru", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(games, nil)
	repo.On("Save", ctx, mock.AnythingOfType("domain.Card")).Return(assert.AnError)

	svc := NewGetDetailedCardService(client, repo, cache)
	_, err := svc.Execute(ctx, "hikaru")

	assert.Error(t, err)
}

func TestGetDetailedCardService_BuildDetailedCard_ReturnsErrorWhenProfileFetchFails(t *testing.T) {
	ctx := context.Background()

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	cache.On("GetCard", ctx, "unknown").Return(domain.Card{}, false, nil)
	repo.On("FindByUsername", ctx, "unknown").Return(domain.Card{}, false, nil)
	client.On("GetProfile", ctx, "unknown").Return(domain.Player{}, assert.AnError)

	svc := NewGetDetailedCardService(client, repo, cache)
	_, err := svc.Execute(ctx, "unknown")

	assert.Error(t, err)
	client.AssertNotCalled(t, "GetStats")
}

func TestGetDetailedCardService_BuildDetailedCard_ReturnsErrorWhenStatsFetchFails(t *testing.T) {
	ctx := context.Background()
	player := domain.Player{Username: "hikaru", Title: domain.TitleGM}

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	cache.On("GetCard", ctx, "hikaru").Return(domain.Card{}, false, nil)
	repo.On("FindByUsername", ctx, "hikaru").Return(domain.Card{}, false, nil)
	client.On("GetProfile", ctx, "hikaru").Return(player, nil)
	client.On("GetStats", ctx, "hikaru").Return(domain.PlayerStats{}, assert.AnError)

	svc := NewGetDetailedCardService(client, repo, cache)
	_, err := svc.Execute(ctx, "hikaru")

	assert.Error(t, err)
	client.AssertNotCalled(t, "GetGames")
}

func TestGetDetailedCardService_PromoteIfPopular_SilentlyIgnoresViewCountError(t *testing.T) {
	ctx := context.Background()
	card := domain.Card{Player: domain.Player{Username: "someuser", Title: domain.TitleNone}}

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	cache.On("IncrementViewCount", ctx, "someuser").Return(0, assert.AnError)

	svc := NewGetDetailedCardService(client, repo, cache)
	svc.promoteIfPopular(ctx, "someuser", card)

	cache.AssertNotCalled(t, "SetCard")
}
