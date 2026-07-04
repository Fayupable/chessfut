package chesscom

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fayupable/chessfut-be/domain"
)

func TestExtractCountryCode(t *testing.T) {
	t.Run("extracts code from country URL", func(t *testing.T) {
		code := extractCountryCode("https://api.chess.com/pub/country/US")

		assert.Equal(t, "US", code)
	})

	t.Run("returns empty string when URL does not match", func(t *testing.T) {
		code := extractCountryCode("")

		assert.Equal(t, "", code)
	})
}

func TestMapResult(t *testing.T) {
	assert.Equal(t, domain.ResultWin, mapResult("win"))
	assert.Equal(t, domain.ResultDraw, mapResult("agreed"))
	assert.Equal(t, domain.ResultDraw, mapResult("repetition"))
	assert.Equal(t, domain.ResultDraw, mapResult("stalemate"))
	assert.Equal(t, domain.ResultDraw, mapResult("insufficient"))
	assert.Equal(t, domain.ResultDraw, mapResult("50move"))
	assert.Equal(t, domain.ResultDraw, mapResult("timevsinsufficient"))
	assert.Equal(t, domain.ResultLoss, mapResult("resigned"))
	assert.Equal(t, domain.ResultLoss, mapResult("checkmated"))
	assert.Equal(t, domain.ResultLoss, mapResult("timeout"))
	assert.Equal(t, domain.ResultLoss, mapResult("abandoned"))
}

func TestExtractOpeningName(t *testing.T) {
	t.Run("strips move notation suffix from slug", func(t *testing.T) {
		url := "https://www.chess.com/openings/Closed-Sicilian-Defense-Grand-Prix-Attack-3...g6-4.Bc4-Bg7-5.Nf3"

		name := extractOpeningName(url)

		assert.Equal(t, "Closed Sicilian Defense Grand Prix Attack", name)
	})

	t.Run("handles URL with no move notation", func(t *testing.T) {
		url := "https://www.chess.com/openings/Italian-Game"

		name := extractOpeningName(url)

		assert.Equal(t, "Italian Game", name)
	})
}

func TestExtractECOCode(t *testing.T) {
	t.Run("extracts ECO code from PGN tag", func(t *testing.T) {
		pgn := `[Event "Live Chess"]
[ECO "B23"]
[Result "1-0"]

1. e4 c5 1-0`

		code := extractECOCode(pgn)

		assert.Equal(t, "B23", code)
	})

	t.Run("returns empty string when ECO tag is missing", func(t *testing.T) {
		code := extractECOCode(`[Event "Live Chess"]`)

		assert.Equal(t, "", code)
	})
}

func TestCountMoves(t *testing.T) {
	t.Run("counts full moves excluding clock annotations and move numbers", func(t *testing.T) {
		pgn := `[Event "Live Chess"]
[ECO "B23"]

1. e4 {[%clk 0:03:00]} 1... c5 {[%clk 0:03:00]} 2. Nc3 {[%clk 0:02:59.9]} 2... g6 {[%clk 0:02:57.5]} 1-0`

		moves := countMoves(pgn)

		assert.Equal(t, 2, moves)
	})

	t.Run("returns zero for empty movetext", func(t *testing.T) {
		moves := countMoves(`[Event "Live Chess"]

*`)

		assert.Equal(t, 0, moves)
	})
}
