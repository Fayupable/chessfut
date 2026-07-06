package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/fayupable/chessfut-be/domain"
)

func TestRefreshCardService_Execute(t *testing.T) {
	ctx := context.Background()
	player := domain.Player{Username: "hikaru", Title: domain.TitleGM}
	stats := domain.PlayerStats{FideRating: 2814, Blitz: domain.TimeControlStats{Rating: 3414}, Rapid: domain.TimeControlStats{Rating: 2839}}

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	client.On("GetProfile", ctx, "hikaru").Return(player, nil)
	client.On("GetStats", ctx, "hikaru").Return(stats, nil)
	repo.On("Save", ctx, mock.AnythingOfType("domain.Card")).Return(nil)
	cache.On("SetCard", ctx, mock.AnythingOfType("domain.Card")).Return(nil)

	svc := NewRefreshCardService(client, repo, cache)
	result, err := svc.Execute(ctx, "hikaru")

	assert.NoError(t, err)
	assert.Equal(t, "hikaru", result.Player.Username)
	assert.Greater(t, result.OVR, 0)
	assert.LessOrEqual(t, result.OVR, 99)
	client.AssertExpectations(t)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestRefreshCardService_Execute_ReturnsErrorWhenChessComFails(t *testing.T) {
	ctx := context.Background()

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	client.On("GetProfile", ctx, "unknown").Return(domain.Player{}, assert.AnError)

	svc := NewRefreshCardService(client, repo, cache)
	_, err := svc.Execute(ctx, "unknown")

	assert.Error(t, err)
	repo.AssertNotCalled(t, "Save")
}

func TestRefreshCardService_Execute_ReturnsErrorWhenSaveFails(t *testing.T) {
	ctx := context.Background()
	player := domain.Player{Username: "hikaru", Title: domain.TitleGM}
	stats := domain.PlayerStats{FideRating: 2814}

	client := new(mockChessComClient)
	repo := new(mockCardRepository)
	cache := new(mockCache)
	client.On("GetProfile", ctx, "hikaru").Return(player, nil)
	client.On("GetStats", ctx, "hikaru").Return(stats, nil)
	repo.On("Save", ctx, mock.AnythingOfType("domain.Card")).Return(assert.AnError)

	svc := NewRefreshCardService(client, repo, cache)
	_, err := svc.Execute(ctx, "hikaru")

	assert.Error(t, err)
	cache.AssertNotCalled(t, "SetCard")
}
