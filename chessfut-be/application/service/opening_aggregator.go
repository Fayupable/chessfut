package service

import (
	"sort"

	"github.com/fayupable/chessfut-be/domain"
)

const topOpeningsLimit = 5

func aggregateTopOpenings(games []domain.Game) []domain.OpeningStat {
	statsByECO := make(map[string]*domain.OpeningStat)

	for _, g := range games {
		s, exists := statsByECO[g.Opening.ECO]
		if !exists {
			s = &domain.OpeningStat{Opening: g.Opening}
			statsByECO[g.Opening.ECO] = s
		}
		s.Count++
		switch g.Result {
		case domain.ResultWin:
			s.Wins++
		case domain.ResultLoss:
			s.Losses++
		case domain.ResultDraw:
			s.Draws++
		}
	}

	result := make([]domain.OpeningStat, 0, len(statsByECO))
	for _, s := range statsByECO {
		result = append(result, *s)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	if len(result) > topOpeningsLimit {
		result = result[:topOpeningsLimit]
	}
	return result
}
