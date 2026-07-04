package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayer_HasFideTitle(t *testing.T) {
	assert.True(t, Player{Title: TitleGM}.HasFideTitle())
	assert.False(t, Player{Title: TitleNone}.HasFideTitle())
}

func TestTimeControlStats_TotalGames(t *testing.T) {
	stats := TimeControlStats{Wins: 10, Losses: 5, Draws: 3}

	assert.Equal(t, 18, stats.TotalGames())
}

func TestTimeControlStats_WinRate(t *testing.T) {
	t.Run("calculates percentage", func(t *testing.T) {
		stats := TimeControlStats{Wins: 5, Losses: 3, Draws: 2}

		assert.Equal(t, 50.0, stats.WinRate())
	})

	t.Run("returns zero when no games played", func(t *testing.T) {
		stats := TimeControlStats{}

		assert.Equal(t, 0.0, stats.WinRate())
	})
}
