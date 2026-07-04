package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToPosition(t *testing.T) {
	assert.Equal(t, PositionForward, MapToPosition(PlayStyleAggressive))
	assert.Equal(t, PositionMidfielder, MapToPosition(PlayStylePositional))
	assert.Equal(t, PositionDefender, MapToPosition(PlayStyleDefensive))
	assert.Equal(t, PositionAllRounder, MapToPosition(PlayStyleBalanced))
}
