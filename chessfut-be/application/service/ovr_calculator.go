package service

import "github.com/fayupable/chessfut-be/domain"

const (
	fideMin = 1000
	fideMax = 2900

	chesscomMin = 100
	chesscomMax = 3500
)

func normalize(value, min, max int) float64 {
	if value <= min {
		return 0
	}
	if value >= max {
		return 100
	}
	return float64(value-min) / float64(max-min) * 100
}

func CalculateOVR(stats domain.PlayerStats) int {
	blitzScore := normalize(stats.Blitz.Rating, chesscomMin, chesscomMax)
	rapidScore := normalize(stats.Rapid.Rating, chesscomMin, chesscomMax)

	if stats.FideRating > 0 {
		fideScore := normalize(stats.FideRating, fideMin, fideMax)
		weighted := fideScore*0.4 + blitzScore*0.3 + rapidScore*0.3
		return int(weighted)
	}

	bulletScore := normalize(stats.Bullet.Rating, chesscomMin, chesscomMax)
	weighted := blitzScore*0.5 + rapidScore*0.3 + bulletScore*0.2
	return int(weighted)
}
