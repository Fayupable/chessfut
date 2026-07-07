package http

import (
	"github.com/fayupable/chessfut-be/domain"
)

func toCardResponse(c domain.Card) cardResponse {
	openings := make([]openingResponse, 0, len(c.TopOpenings))
	for _, o := range c.TopOpenings {
		openings = append(openings, openingResponse{
			ECO:     o.Opening.ECO,
			Name:    o.Opening.Name,
			Count:   o.Count,
			WinRate: o.WinRate(),
		})
	}

	badges := make([]string, 0, len(c.Badges))
	for _, b := range c.Badges {
		badges = append(badges, string(b))
	}

	return cardResponse{
		Username:            c.Player.Username,
		Name:                c.Player.Name,
		Title:               string(c.Player.Title),
		Avatar:              c.Player.Avatar,
		Followers:           c.Player.Followers,
		CountryCode:         c.Player.CountryCode,
		JoinedAt:            c.Player.JoinedAt,
		ChesscomURL:         c.Player.URL,
		LastOnline:          c.Player.LastOnline,
		CardType:            string(c.CardType),
		Tier:                string(c.Tier),
		OVR:                 c.OVR,
		PlayStyle:           string(c.PlayStyle),
		Position:            string(c.Position),
		TacticsRating:       c.Stats.TacticsRating,
		PuzzleRushAccuracy:  c.Stats.PuzzleRushAccuracy,
		FideRating:          c.Stats.FideRating,
		EffectiveFideRating: c.EffectiveFideRating,
		FideSource:          c.FideSource,
		GamesSnapshot:       c.GamesSnapshot,
		Attributes: attributesResponse{
			Pac: c.Attributes.Pac,
			Sho: c.Attributes.Sho,
			Pas: c.Attributes.Pas,
			Dri: c.Attributes.Dri,
			Def: c.Attributes.Def,
			Phy: c.Attributes.Phy,
		},
		WorkRate: workRateResponse{
			Attack:  string(c.WorkRate.Attack),
			Defense: string(c.WorkRate.Defense),
		},
		Badges:      badges,
		Bullet:      toTimeControlResponse(c.Stats.Bullet),
		Blitz:       toTimeControlResponse(c.Stats.Blitz),
		Rapid:       toTimeControlResponse(c.Stats.Rapid),
		Daily:       toTimeControlResponse(c.Stats.Daily),
		TopOpenings: openings,
		LastUpdated: c.ComputedAt,
	}
}

func toTimeControlResponse(s domain.TimeControlStats) timeControlResponse {
	return timeControlResponse{
		Rating:  s.Rating,
		Highest: s.Highest,
		Wins:    s.Wins,
		Losses:  s.Losses,
		Draws:   s.Draws,
		WinRate: s.WinRate(),
		RD:      s.RD,
	}
}
