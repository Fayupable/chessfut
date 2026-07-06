package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	t.Run("value at or below min returns zero", func(t *testing.T) {
		assert.Equal(t, 0.0, normalize(50, 3500))
	})

	t.Run("value at or above max returns one hundred", func(t *testing.T) {
		assert.Equal(t, 100.0, normalize(4000, 3500))
	})

	t.Run("value in range scales linearly", func(t *testing.T) {
		assert.Equal(t, 50.0, normalize(1800, 3500))
	})
}
