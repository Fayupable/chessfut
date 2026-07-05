package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/fayupable/chessfut-be/domain"
)

func TestGetFastCardService_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("returns card from cache without touching repository or chess.com", func(t *testing.T) {
		cachedCard := domain.Card{Player: domain.Player{Username: "hikaru"}}

		client := new(mockChessComClient)
		repo := new(mockCardRepository)
		cache := new(mockCache)
		cache.On("GetCard", ctx, "hikaru").Return(cachedCard, true, nil)

		svc := NewGetFastCardService(client, repo, cache)
		result, err := svc.Execute(ctx, "hikaru")

		assert.NoError(t, err)
		assert.Equal(t, cachedCard, result)
		cache.AssertExpectations(t)
		repo.AssertNotCalled(t, "FindByUsername")
		client.AssertNotCalled(t, "GetProfile")
	})

	t.Run("returns card from repository and promotes non-titled player below threshold", func(t *testing.T) {
		existingCard := domain.Card{Player: domain.Player{Username: "someuser", Title: domain.TitleNone}}

		client := new(mockChessComClient)
		repo := new(mockCardRepository)
		cache := new(mockCache)
		cache.On("GetCard", ctx, "someuser").Return(domain.Card{}, false, nil)
		repo.On("FindByUsername", ctx, "someuser").Return(existingCard, true, nil)
		cache.On("IncrementViewCount", ctx, "someuser").Return(3, nil)

		svc := NewGetFastCardService(client, repo, cache)
		result, err := svc.Execute(ctx, "someuser")

		assert.NoError(t, err)
		assert.Equal(t, existingCard, result)
		cache.AssertExpectations(t)
		cache.AssertNotCalled(t, "SetCard")
	})

	t.Run("builds card from chess.com and immediately caches titled player", func(t *testing.T) {
		player := domain.Player{Username: "hikaru", Title: domain.TitleGM}
		stats := domain.PlayerStats{FideRating: 2814, Blitz: domain.TimeControlStats{Rating: 3414}, Rapid: domain.TimeControlStats{Rating: 2839}}

		client := new(mockChessComClient)
		repo := new(mockCardRepository)
		cache := new(mockCache)
		cache.On("GetCard", ctx, "hikaru").Return(domain.Card{}, false, nil)
		repo.On("FindByUsername", ctx, "hikaru").Return(domain.Card{}, false, nil)
		client.On("GetProfile", ctx, "hikaru").Return(player, nil)
		client.On("GetStats", ctx, "hikaru").Return(stats, nil)
		repo.On("Save", ctx, mock.AnythingOfType("domain.Card")).Return(nil)
		cache.On("SetCard", ctx, mock.AnythingOfType("domain.Card")).Return(nil)

		svc := NewGetFastCardService(client, repo, cache)
		result, err := svc.Execute(ctx, "hikaru")

		assert.NoError(t, err)
		assert.Equal(t, "hikaru", result.Player.Username)
		assert.Equal(t, domain.CardTierTitled, result.Tier)
		assert.Greater(t, result.OVR, 0)
		assert.LessOrEqual(t, result.OVR, 99)
		client.AssertExpectations(t)
		repo.AssertExpectations(t)
		cache.AssertExpectations(t)
	})
}

func TestGetFastCardService_Execute_PromotesNonTitledPlayerAtThreshold(t *testing.T) {
	ctx := context.Background()
	player := domain.Player{Username: "someuser", Title: domain.TitleNone}
	stats := domain.PlayerStats{Blitz: domain.TimeControlStats{Rating: 1500}}

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	cache.On("GetCard", ctx, "someuser").Return(domain.Card{}, false, nil)
	repo.On("FindByUsername", ctx, "someuser").Return(domain.Card{}, false, nil)
	client.On("GetProfile", ctx, "someuser").Return(player, nil)
	client.On("GetStats", ctx, "someuser").Return(stats, nil)
	repo.On("Save", ctx, mock.AnythingOfType("domain.Card")).Return(nil)
	cache.On("IncrementViewCount", ctx, "someuser").Return(5, nil)
	cache.On("SetCard", ctx, mock.AnythingOfType("domain.Card")).Return(nil)

	svc := NewGetFastCardService(client, repo, cache)
	_, err := svc.Execute(ctx, "someuser")

	assert.NoError(t, err)
	cache.AssertExpectations(t)
}

func TestGetFastCardService_Execute_ReturnsErrorWhenProfileFetchFails(t *testing.T) {
	ctx := context.Background()

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	cache.On("GetCard", ctx, "unknown").Return(domain.Card{}, false, nil)
	repo.On("FindByUsername", ctx, "unknown").Return(domain.Card{}, false, nil)
	client.On("GetProfile", ctx, "unknown").Return(domain.Player{}, assert.AnError)

	svc := NewGetFastCardService(client, repo, cache)
	_, err := svc.Execute(ctx, "unknown")

	assert.Error(t, err)
	repo.AssertNotCalled(t, "Save")
}

func TestGetFastCardService_Execute_ReturnsErrorWhenRepositoryFails(t *testing.T) {
	ctx := context.Background()

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	cache.On("GetCard", ctx, "hikaru").Return(domain.Card{}, false, nil)
	repo.On("FindByUsername", ctx, "hikaru").Return(domain.Card{}, false, assert.AnError)

	svc := NewGetFastCardService(client, repo, cache)
	_, err := svc.Execute(ctx, "hikaru")

	assert.Error(t, err)
	client.AssertNotCalled(t, "GetProfile")
}

func TestGetFastCardService_Execute_ReturnsErrorWhenStatsFetchFails(t *testing.T) {
	ctx := context.Background()
	player := domain.Player{Username: "hikaru", Title: domain.TitleGM}

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	cache.On("GetCard", ctx, "hikaru").Return(domain.Card{}, false, nil)
	repo.On("FindByUsername", ctx, "hikaru").Return(domain.Card{}, false, nil)
	client.On("GetProfile", ctx, "hikaru").Return(player, nil)
	client.On("GetStats", ctx, "hikaru").Return(domain.PlayerStats{}, assert.AnError)

	svc := NewGetFastCardService(client, repo, cache)
	_, err := svc.Execute(ctx, "hikaru")

	assert.Error(t, err)
	repo.AssertNotCalled(t, "Save")
}

func TestGetFastCardService_PromoteIfPopular_SilentlyIgnoresViewCountError(t *testing.T) {
	ctx := context.Background()
	card := domain.Card{Player: domain.Player{Username: "someuser", Title: domain.TitleNone}}

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	cache.On("IncrementViewCount", ctx, "someuser").Return(0, assert.AnError)

	svc := NewGetFastCardService(client, repo, cache)
	svc.promoteIfPopular(ctx, "someuser", card)

	cache.AssertNotCalled(t, "SetCard")
}
