package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fayupable/chessfut-be/domain"
)

func TestCalculateOVR(t *testing.T) {
	t.Run("titled player uses FIDE-weighted formula", func(t *testing.T) {
		stats := domain.PlayerStats{
			FideRating: 2814,
			Blitz:      domain.TimeControlStats{Rating: 3414},
			Rapid:      domain.TimeControlStats{Rating: 2839},
		}

		ovr := CalculateOVR(stats)

		assert.Equal(t, 91, ovr)
	})

	t.Run("non-titled player uses chess.com-only formula", func(t *testing.T) {
		stats := domain.PlayerStats{
			FideRating: 0,
			Blitz:      domain.TimeControlStats{Rating: 2000},
			Rapid:      domain.TimeControlStats{Rating: 1500},
			Bullet:     domain.TimeControlStats{Rating: 1000},
		}

		ovr := CalculateOVR(stats)

		assert.Equal(t, 45, ovr)
	})

	t.Run("rating below normalization floor clamps to zero contribution", func(t *testing.T) {
		stats := domain.PlayerStats{
			FideRating: 0,
			Blitz:      domain.TimeControlStats{Rating: 50},
			Rapid:      domain.TimeControlStats{Rating: 50},
			Bullet:     domain.TimeControlStats{Rating: 50},
		}

		ovr := CalculateOVR(stats)

		assert.Equal(t, 0, ovr)
	})

	t.Run("rating above normalization ceiling clamps to full contribution", func(t *testing.T) {
		stats := domain.PlayerStats{
			FideRating: 0,
			Blitz:      domain.TimeControlStats{Rating: 4000},
			Rapid:      domain.TimeControlStats{Rating: 4000},
			Bullet:     domain.TimeControlStats{Rating: 4000},
		}

		ovr := CalculateOVR(stats)

		assert.Equal(t, 100, ovr)
	})
}
