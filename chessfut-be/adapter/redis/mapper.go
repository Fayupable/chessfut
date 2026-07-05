package redis

import (
	"github.com/fayupable/chessfut-be/domain"
)

func toCardModel(c domain.Card) cardModel {
	badges := make([]string, 0, len(c.Badges))
	for _, b := range c.Badges {
		badges = append(badges, string(b))
	}

	openings := make([]openingStatModel, 0, len(c.TopOpenings))
	for _, o := range c.TopOpenings {
		openings = append(openings, openingStatModel{
			ECO: o.Opening.ECO, Name: o.Opening.Name,
			Count: o.Count, Wins: o.Wins, Losses: o.Losses, Draws: o.Draws,
		})
	}

	return cardModel{
		Username:           c.Player.Username,
		Name:               c.Player.Name,
		Title:              string(c.Player.Title),
		Avatar:             c.Player.Avatar,
		Followers:          c.Player.Followers,
		CountryCode:        c.Player.CountryCode,
		JoinedAt:           c.Player.JoinedAt,
		FideRating:         c.Stats.FideRating,
		TacticsRating:      c.Stats.TacticsRating,
		PuzzleRushAccuracy: c.Stats.PuzzleRushAccuracy,
		Bullet:             toTimeControlModel(c.Stats.Bullet),
		Blitz:              toTimeControlModel(c.Stats.Blitz),
		Rapid:              toTimeControlModel(c.Stats.Rapid),
		Daily:              toTimeControlModel(c.Stats.Daily),
		CardType:           string(c.CardType),
		Tier:               string(c.Tier),
		OVR:                c.OVR,
		PlayStyle:          string(c.PlayStyle),
		Position:           string(c.Position),
		Attributes: attributesModel{
			Pac: c.Attributes.Pac,
			Sho: c.Attributes.Sho,
			Pas: c.Attributes.Pas,
			Dri: c.Attributes.Dri,
			Def: c.Attributes.Def,
			Phy: c.Attributes.Phy,
		},
		WorkRate: workRateModel{
			Attack:  string(c.WorkRate.Attack),
			Defense: string(c.WorkRate.Defense),
		},
		Badges:        badges,
		TopOpenings:   openings,
		GamesSnapshot: c.GamesSnapshot,
		ComputedAt:    c.ComputedAt,
		ExpiresAt:     c.ExpiresAt,
	}
}

func toTimeControlModel(s domain.TimeControlStats) timeControlStatsModel {
	return timeControlStatsModel{Rating: s.Rating, Highest: s.Highest, Wins: s.Wins, Losses: s.Losses, Draws: s.Draws}
}

func fromCardModel(m cardModel) domain.Card {
	badges := make([]domain.Badge, 0, len(m.Badges))
	for _, b := range m.Badges {
		badges = append(badges, domain.Badge(b))
	}

	openings := make([]domain.OpeningStat, 0, len(m.TopOpenings))
	for _, o := range m.TopOpenings {
		openings = append(openings, domain.OpeningStat{
			Opening: domain.Opening{ECO: o.ECO, Name: o.Name},
			Count:   o.Count, Wins: o.Wins, Losses: o.Losses, Draws: o.Draws,
		})
	}

	return domain.Card{
		Player: domain.Player{
			Username: m.Username, Name: m.Name, Title: domain.Title(m.Title),
			Avatar: m.Avatar, Followers: m.Followers, CountryCode: m.CountryCode, JoinedAt: m.JoinedAt,
		},
		Stats: domain.PlayerStats{
			FideRating:         m.FideRating,
			TacticsRating:      m.TacticsRating,
			PuzzleRushAccuracy: m.PuzzleRushAccuracy,
			Bullet:             fromTimeControlModel(m.Bullet, domain.TimeControlBullet),
			Blitz:              fromTimeControlModel(m.Blitz, domain.TimeControlBlitz),
			Rapid:              fromTimeControlModel(m.Rapid, domain.TimeControlRapid),
			Daily:              fromTimeControlModel(m.Daily, domain.TimeControlDaily),
		},
		CardType:  domain.CardType(m.CardType),
		Tier:      domain.CardTier(m.Tier),
		OVR:       m.OVR,
		PlayStyle: domain.PlayStyle(m.PlayStyle),
		Position:  domain.Position(m.Position),
		Attributes: domain.Attributes{
			Pac: m.Attributes.Pac,
			Sho: m.Attributes.Sho,
			Pas: m.Attributes.Pas,
			Dri: m.Attributes.Dri,
			Def: m.Attributes.Def,
			Phy: m.Attributes.Phy,
		},
		WorkRate: domain.WorkRate{
			Attack:  domain.WorkRateLevel(m.WorkRate.Attack),
			Defense: domain.WorkRateLevel(m.WorkRate.Defense),
		},
		Badges:        badges,
		TopOpenings:   openings,
		GamesSnapshot: m.GamesSnapshot,
		ComputedAt:    m.ComputedAt,
		ExpiresAt:     m.ExpiresAt,
	}
}

func fromTimeControlModel(m timeControlStatsModel, tc domain.TimeControl) domain.TimeControlStats {
	return domain.TimeControlStats{TimeControl: tc, Rating: m.Rating, Highest: m.Highest, Wins: m.Wins, Losses: m.Losses, Draws: m.Draws}
}
