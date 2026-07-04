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
		Username:    c.Player.Username,
		Name:        c.Player.Name,
		Title:       string(c.Player.Title),
		Avatar:      c.Player.Avatar,
		Followers:   c.Player.Followers,
		CountryCode: c.Player.CountryCode,
		JoinedAt:    c.Player.JoinedAt,
		CardType:    string(c.CardType),
		Tier:        string(c.Tier),
		OVR:         c.OVR,
		PlayStyle:   string(c.PlayStyle),
		Position:    string(c.Position),
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
	}
}
