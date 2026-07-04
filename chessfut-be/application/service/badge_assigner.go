package service

import "github.com/fayupable/chessfut-be/domain"

func AssignBadges(player domain.Player, stats domain.PlayerStats) []domain.Badge {
	var badges []domain.Badge

	if player.HasFideTitle() {
		badges = append(badges, domain.BadgeTitled)
	}

	if stats.Blitz.Rating >= 2800 {
		badges = append(badges, domain.BadgeBlitzLegend)
	}

	if stats.Bullet.Rating >= 3000 {
		badges = append(badges, domain.BadgeBulletBeast)
	}

	if stats.Daily.TotalGames() >= 500 {
		badges = append(badges, domain.BadgeMarathoner)
	}

	if stats.FideRating > 0 && stats.Blitz.Rating > stats.FideRating+400 {
		badges = append(badges, domain.BadgeGiantSlayer)
	}

	return badges
}
