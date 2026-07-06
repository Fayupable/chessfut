package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fayupable/chessfut-be/domain"
)

func TestGetLeaderboardService_Execute(t *testing.T) {
	ctx := context.Background()
	expected := []domain.Card{
		{Player: domain.Player{Username: "hikaru"}, OVR: 91},
		{Player: domain.Player{Username: "magnuscarlsen"}, OVR: 89},
	}

	repo := new(mockCardRepository)
	repo.On("FindTopByOVR", ctx, 50, 0).Return(expected, nil)

	svc := NewGetLeaderboardService(repo)
	result, err := svc.Execute(ctx, 50, 0)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}
