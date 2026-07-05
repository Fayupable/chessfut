package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fayupable/chessfut-be/domain"
)

func TestCalculateWorkRate(t *testing.T) {
	t.Run("high attack and defense", func(t *testing.T) {
		attrs := domain.Attributes{Pac: 80, Sho: 80, Def: 80}
		wr := CalculateWorkRate(attrs)

		assert.Equal(t, domain.WorkRateHigh, wr.Attack)
		assert.Equal(t, domain.WorkRateHigh, wr.Defense)
	})

	t.Run("medium attack and defense", func(t *testing.T) {
		attrs := domain.Attributes{Pac: 55, Sho: 55, Def: 55}
		wr := CalculateWorkRate(attrs)

		assert.Equal(t, domain.WorkRateMed, wr.Attack)
		assert.Equal(t, domain.WorkRateMed, wr.Defense)
	})

	t.Run("low attack and defense", func(t *testing.T) {
		attrs := domain.Attributes{Pac: 20, Sho: 20, Def: 20}
		wr := CalculateWorkRate(attrs)

		assert.Equal(t, domain.WorkRateLow, wr.Attack)
		assert.Equal(t, domain.WorkRateLow, wr.Defense)
	})
}

func TestOpeningDiversity(t *testing.T) {
	t.Run("empty openings returns zero", func(t *testing.T) {
		assert.Equal(t, 0.0, openingDiversity(nil))
	})

	t.Run("single opening has low diversity", func(t *testing.T) {
		openings := []domain.OpeningStat{{Opening: domain.Opening{ECO: "A01"}, Count: 100}}
		diversity := openingDiversity(openings)

		assert.Equal(t, 0.0, diversity)
	})

	t.Run("five evenly split openings have high diversity", func(t *testing.T) {
		openings := []domain.OpeningStat{
			{Opening: domain.Opening{ECO: "A01"}, Count: 20},
			{Opening: domain.Opening{ECO: "B01"}, Count: 20},
			{Opening: domain.Opening{ECO: "C01"}, Count: 20},
			{Opening: domain.Opening{ECO: "D01"}, Count: 20},
			{Opening: domain.Opening{ECO: "E01"}, Count: 20},
		}
		diversity := openingDiversity(openings)

		assert.InDelta(t, 99, diversity, 1)
	})
}

func TestClampFloat(t *testing.T) {
	assert.Equal(t, 0.0, clampFloat(-5, 0, 99))
	assert.Equal(t, 99.0, clampFloat(150, 0, 99))
	assert.Equal(t, 50.0, clampFloat(50, 0, 99))
}
