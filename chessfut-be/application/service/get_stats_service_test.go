package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStatsService_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("returns cached count when available", func(t *testing.T) {
		repo := new(mockCardRepository)
		cache := new(mockCache)
		cache.On("GetTotalCardsCount", ctx).Return(42, true, nil)

		svc := NewGetStatsService(repo, cache)
		count, err := svc.Execute(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 42, count)
		repo.AssertNotCalled(t, "CountAll")
	})

	t.Run("falls back to repository and caches result on cache miss", func(t *testing.T) {
		repo := new(mockCardRepository)
		cache := new(mockCache)
		cache.On("GetTotalCardsCount", ctx).Return(0, false, nil)
		repo.On("CountAll", ctx).Return(100, nil)
		cache.On("SetTotalCardsCount", ctx, 100).Return(nil)

		svc := NewGetStatsService(repo, cache)
		count, err := svc.Execute(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 100, count)
	})
}
