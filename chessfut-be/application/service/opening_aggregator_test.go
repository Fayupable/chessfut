package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fayupable/chessfut-be/domain"
)

func gameWithOpening(eco string, result domain.GameResult) domain.Game {
	return domain.Game{Opening: domain.Opening{ECO: eco, Name: eco + "-name"}, Result: result}
}

func TestAggregateTopOpenings(t *testing.T) {
	t.Run("aggregates count and results per ECO code", func(t *testing.T) {
		games := []domain.Game{
			gameWithOpening("B23", domain.ResultWin),
			gameWithOpening("B23", domain.ResultWin),
			gameWithOpening("B23", domain.ResultLoss),
			gameWithOpening("C50", domain.ResultDraw),
		}

		result := aggregateTopOpenings(games)

		assert.Len(t, result, 2)
		assert.Equal(t, "B23", result[0].Opening.ECO)
		assert.Equal(t, 3, result[0].Count)
		assert.Equal(t, 2, result[0].Wins)
		assert.Equal(t, 1, result[0].Losses)
	})

	t.Run("sorts by count descending", func(t *testing.T) {
		games := append(
			[]domain.Game{gameWithOpening("A01", domain.ResultWin)},
			append(
				gamesOfOpening("B23", 5),
				gamesOfOpening("C50", 3)...,
			)...,
		)

		result := aggregateTopOpenings(games)

		assert.Equal(t, "B23", result[0].Opening.ECO)
		assert.Equal(t, "C50", result[1].Opening.ECO)
		assert.Equal(t, "A01", result[2].Opening.ECO)
	})

	t.Run("limits result to top 5 openings", func(t *testing.T) {
		var games []domain.Game
		for i := 0; i < 7; i++ {
			games = append(games, gamesOfOpening(string(rune('A'+i)), i+1)...)
		}

		result := aggregateTopOpenings(games)

		assert.Len(t, result, 5)
	})
}

func gamesOfOpening(eco string, count int) []domain.Game {
	games := make([]domain.Game, count)
	for i := range games {
		games[i] = gameWithOpening(eco, domain.ResultWin)
	}
	return games
}
