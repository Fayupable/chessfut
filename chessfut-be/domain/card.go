package domain

import "time"

type CardTier string

const (
	CardTierStandard CardTier = "standard"
	CardTierTitled   CardTier = "titled"
)

type CardType string

const (
	CardTypeFast     CardType = "fast"
	CardTypeDetailed CardType = "detailed"
)

type PlayStyle string

const (
	PlayStyleAggressive PlayStyle = "aggressive"
	PlayStylePositional PlayStyle = "positional"
	PlayStyleDefensive  PlayStyle = "defensive"
	PlayStyleBalanced   PlayStyle = "balanced"
)

type Badge string

const (
	BadgeBlitzLegend Badge = "blitz_legend"
	BadgeTitled      Badge = "titled"
	BadgeBulletBeast Badge = "bullet_beast"
	BadgeMarathoner  Badge = "marathoner"
	BadgeGiantSlayer Badge = "giant_slayer"
)

type Card struct {
	Player        Player
	Stats         PlayerStats
	CardType      CardType
	Tier          CardTier
	OVR           int
	PlayStyle     PlayStyle
	Position      Position
	Badges        []Badge
	TopOpenings   []OpeningStat
	GamesSnapshot int
	ComputedAt    time.Time
	ExpiresAt     time.Time
}

func TotalGames(stats PlayerStats) int {
	return stats.Bullet.TotalGames() + stats.Blitz.TotalGames() + stats.Rapid.TotalGames() + stats.Daily.TotalGames()
}

func (c Card) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}
