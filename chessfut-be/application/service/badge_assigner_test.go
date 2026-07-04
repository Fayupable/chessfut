package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fayupable/chessfut-be/domain"
)

func TestAssignBadges(t *testing.T) {
	t.Run("untitled player with low stats gets no badges", func(t *testing.T) {
		player := domain.Player{Title: domain.TitleNone}
		stats := domain.PlayerStats{}

		badges := AssignBadges(player, stats)

		assert.Empty(t, badges)
	})

	t.Run("titled player receives titled badge", func(t *testing.T) {
		player := domain.Player{Title: domain.TitleGM}
		stats := domain.PlayerStats{}

		badges := AssignBadges(player, stats)

		assert.Contains(t, badges, domain.BadgeTitled)
	})

	t.Run("high blitz rating receives blitz legend badge", func(t *testing.T) {
		player := domain.Player{}
		stats := domain.PlayerStats{Blitz: domain.TimeControlStats{Rating: 2900}}

		badges := AssignBadges(player, stats)

		assert.Contains(t, badges, domain.BadgeBlitzLegend)
	})

	t.Run("high bullet rating receives bullet beast badge", func(t *testing.T) {
		player := domain.Player{}
		stats := domain.PlayerStats{Bullet: domain.TimeControlStats{Rating: 3100}}

		badges := AssignBadges(player, stats)

		assert.Contains(t, badges, domain.BadgeBulletBeast)
	})

	t.Run("500+ daily games receives marathoner badge", func(t *testing.T) {
		player := domain.Player{}
		stats := domain.PlayerStats{Daily: domain.TimeControlStats{Wins: 300, Losses: 150, Draws: 60}}

		badges := AssignBadges(player, stats)

		assert.Contains(t, badges, domain.BadgeMarathoner)
	})

	t.Run("blitz rating far above FIDE rating receives giant slayer badge", func(t *testing.T) {
		player := domain.Player{}
		stats := domain.PlayerStats{FideRating: 2000, Blitz: domain.TimeControlStats{Rating: 2500}}

		badges := AssignBadges(player, stats)

		assert.Contains(t, badges, domain.BadgeGiantSlayer)
	})
}
