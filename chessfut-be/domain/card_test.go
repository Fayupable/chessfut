package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTotalGames(t *testing.T) {
	stats := PlayerStats{
		Bullet: TimeControlStats{Wins: 1, Losses: 1, Draws: 1},
		Blitz:  TimeControlStats{Wins: 2, Losses: 2, Draws: 2},
		Rapid:  TimeControlStats{Wins: 3, Losses: 3, Draws: 3},
		Daily:  TimeControlStats{Wins: 4, Losses: 4, Draws: 4},
	}

	assert.Equal(t, 30, TotalGames(stats))
}

func TestCard_IsExpired(t *testing.T) {
	t.Run("expired when expires_at is in the past", func(t *testing.T) {
		card := Card{ExpiresAt: time.Now().Add(-1 * time.Hour)}

		assert.True(t, card.IsExpired())
	})

	t.Run("not expired when expires_at is in the future", func(t *testing.T) {
		card := Card{ExpiresAt: time.Now().Add(1 * time.Hour)}

		assert.False(t, card.IsExpired())
	})
}
