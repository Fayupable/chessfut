package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/fayupable/chessfut-be/domain"
)

func TestSyncTitledPlayersService_Execute(t *testing.T) {
	ctx := context.Background()

	client := new(mockChessComClient)
	repo := new(mockCardRepository)

	for _, title := range syncedTitles {
		if title == domain.TitleGM {
			client.On("GetTitledUsernames", ctx, title).Return([]string{"hikaru"}, nil)
			continue
		}
		client.On("GetTitledUsernames", ctx, title).Return([]string{}, nil)
	}

	player := domain.Player{Username: "hikaru", Title: domain.TitleGM}
	stats := domain.PlayerStats{FideRating: 2814, Blitz: domain.TimeControlStats{Rating: 3414}, Rapid: domain.TimeControlStats{Rating: 2839}}

	repo.On("FindByUsername", ctx, "hikaru").Return(domain.Card{}, false, nil)
	client.On("GetProfile", ctx, "hikaru").Return(player, nil)
	client.On("GetStats", ctx, "hikaru").Return(stats, nil)
	repo.On("Save", ctx, mock.AnythingOfType("domain.Card")).Return(nil)

	svc := NewSyncTitledPlayersService(client, repo)
	synced, err := svc.Execute(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 1, synced)
	client.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestSyncTitledPlayersService_Execute_SkipsTitleWhenFetchFails(t *testing.T) {
	ctx := context.Background()

	client := new(mockChessComClient)
	repo := new(mockCardRepository)

	for _, title := range syncedTitles {
		if title == domain.TitleGM {
			client.On("GetTitledUsernames", ctx, title).Return([]string(nil), assert.AnError)
			continue
		}
		client.On("GetTitledUsernames", ctx, title).Return([]string{}, nil)
	}

	svc := NewSyncTitledPlayersService(client, repo)
	synced, err := svc.Execute(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 0, synced)
	repo.AssertNotCalled(t, "FindByUsername")
}

func TestSyncTitledPlayersService_Execute_SkipsUsernameAlreadyInRepository(t *testing.T) {
	ctx := context.Background()

	client := new(mockChessComClient)
	repo := new(mockCardRepository)

	for _, title := range syncedTitles {
		if title == domain.TitleGM {
			client.On("GetTitledUsernames", ctx, title).Return([]string{"hikaru"}, nil)
			continue
		}
		client.On("GetTitledUsernames", ctx, title).Return([]string{}, nil)
	}
	repo.On("FindByUsername", ctx, "hikaru").Return(domain.Card{}, true, nil)

	svc := NewSyncTitledPlayersService(client, repo)
	synced, err := svc.Execute(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 0, synced)
	client.AssertNotCalled(t, "GetProfile")
}

func TestSyncTitledPlayersService_Execute_SkipsUsernameWhenSaveFails(t *testing.T) {
	ctx := context.Background()
	player := domain.Player{Username: "hikaru", Title: domain.TitleGM}
	stats := domain.PlayerStats{FideRating: 2814}

	client := new(mockChessComClient)
	repo := new(mockCardRepository)

	for _, title := range syncedTitles {
		if title == domain.TitleGM {
			client.On("GetTitledUsernames", ctx, title).Return([]string{"hikaru"}, nil)
			continue
		}
		client.On("GetTitledUsernames", ctx, title).Return([]string{}, nil)
	}
	repo.On("FindByUsername", ctx, "hikaru").Return(domain.Card{}, false, nil)
	client.On("GetProfile", ctx, "hikaru").Return(player, nil)
	client.On("GetStats", ctx, "hikaru").Return(stats, nil)
	repo.On("Save", ctx, mock.AnythingOfType("domain.Card")).Return(assert.AnError)

	svc := NewSyncTitledPlayersService(client, repo)
	synced, err := svc.Execute(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 0, synced)
}
