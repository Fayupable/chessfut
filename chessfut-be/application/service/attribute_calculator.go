package service

import "github.com/fayupable/chessfut-be/domain"

func CalculateWorkRate(attrs domain.Attributes) domain.WorkRate {
	attack := workRateLevel((attrs.Pac + attrs.Sho) / 2)
	defense := workRateLevel(attrs.Def)
	return domain.WorkRate{Attack: attack, Defense: defense}
}

func workRateLevel(score int) domain.WorkRateLevel {
	switch {
	case score >= 68:
		return domain.WorkRateHigh
	case score >= 50:
		return domain.WorkRateMed
	default:
		return domain.WorkRateLow
	}
}

func openingDiversity(topOpenings []domain.OpeningStat) float64 {
	if len(topOpenings) == 0 {
		return 0
	}

	total := 0
	for _, o := range topOpenings {
		total += o.Count
	}
	if total == 0 {
		return 0
	}

	sumSquares := 0.0
	for _, o := range topOpenings {
		share := float64(o.Count) / float64(total)
		sumSquares += share * share
	}

	simpsonIndex := 1 - sumSquares
	breadthBonus := float64(len(topOpenings)) / float64(topOpeningsLimit)

	return clampFloat(simpsonIndex*breadthBonus*99*1.25, 0, 99)
}
