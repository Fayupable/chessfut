package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpeningStat_WinRate(t *testing.T) {
	t.Run("calculates percentage", func(t *testing.T) {
		stat := OpeningStat{Count: 4, Wins: 3}

		assert.Equal(t, 75.0, stat.WinRate())
	})

	t.Run("returns zero when count is zero", func(t *testing.T) {
		stat := OpeningStat{}

		assert.Equal(t, 0.0, stat.WinRate())
	})
}
