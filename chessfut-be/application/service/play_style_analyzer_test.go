package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fayupable/chessfut-be/domain"
)

func gamesWith(count int, movesCount int, result domain.GameResult) []domain.Game {
	games := make([]domain.Game, count)
	for i := range games {
		games[i] = domain.Game{MovesCount: movesCount, Result: result}
	}
	return games
}

func TestDeterminePlayStyle(t *testing.T) {
	t.Run("no games returns balanced", func(t *testing.T) {
		style := DeterminePlayStyle(nil)

		assert.Equal(t, domain.PlayStyleBalanced, style)
	})

	t.Run("high draw rate returns defensive regardless of move count", func(t *testing.T) {
		games := append(gamesWith(2, 20, domain.ResultDraw), gamesWith(2, 20, domain.ResultWin)...)

		style := DeterminePlayStyle(games)

		assert.Equal(t, domain.PlayStyleDefensive, style)
	})

	t.Run("short average game length with low draw rate returns aggressive", func(t *testing.T) {
		games := gamesWith(10, 20, domain.ResultWin)

		style := DeterminePlayStyle(games)

		assert.Equal(t, domain.PlayStyleAggressive, style)
	})

	t.Run("long average game length with low draw rate returns positional", func(t *testing.T) {
		games := gamesWith(10, 50, domain.ResultWin)

		style := DeterminePlayStyle(games)

		assert.Equal(t, domain.PlayStylePositional, style)
	})

	t.Run("moderate average game length with low draw rate returns balanced", func(t *testing.T) {
		games := gamesWith(10, 35, domain.ResultWin)

		style := DeterminePlayStyle(games)

		assert.Equal(t, domain.PlayStyleBalanced, style)
	})
}
