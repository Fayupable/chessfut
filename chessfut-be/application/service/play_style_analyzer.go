package service

import "github.com/fayupable/chessfut-be/domain"

func DeterminePlayStyle(games []domain.Game) domain.PlayStyle {
	if len(games) == 0 {
		return domain.PlayStyleBalanced
	}

	var totalMoves, decisive, draws int
	for _, g := range games {
		totalMoves += g.MovesCount
		if g.Result == domain.ResultDraw {
			draws++
		} else {
			decisive++
		}
	}

	avgMoves := float64(totalMoves) / float64(len(games))
	drawRate := float64(draws) / float64(len(games)) * 100

	switch {
	case drawRate > 25:
		return domain.PlayStyleDefensive
	case avgMoves < 30:
		return domain.PlayStyleAggressive
	case avgMoves > 45:
		return domain.PlayStylePositional
	default:
		return domain.PlayStyleBalanced
	}
}
