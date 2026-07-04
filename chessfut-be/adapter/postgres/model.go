package postgres

import (
	"time"

	"github.com/fayupable/chessfut-be/domain"
)

type openingStatModel struct {
	ECO    string `json:"eco"`
	Name   string `json:"name"`
	Count  int    `json:"count"`
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
	Draws  int    `json:"draws"`
}

type timeControlStatsModel struct {
	Rating  int `json:"rating"`
	Highest int `json:"highest"`
	Wins    int `json:"wins"`
	Losses  int `json:"losses"`
	Draws   int `json:"draws"`
}

type cardDataModel struct {
	Username    string                `json:"username"`
	Name        string                `json:"name"`
	Title       string                `json:"title"`
	Avatar      string                `json:"avatar"`
	Followers   int                   `json:"followers"`
	CountryCode string                `json:"country_code"`
	JoinedAt    time.Time             `json:"joined_at"`
	FideRating  int                   `json:"fide_rating"`
	Bullet      timeControlStatsModel `json:"bullet"`
	Blitz       timeControlStatsModel `json:"blitz"`
	Rapid       timeControlStatsModel `json:"rapid"`
	Daily       timeControlStatsModel `json:"daily"`
	OVR         int                   `json:"ovr"`
	PlayStyle   string                `json:"play_style"`
	Badges      []string              `json:"badges"`
	TopOpenings []openingStatModel    `json:"top_openings"`
}

func toCardDataModel(c domain.Card) cardDataModel {
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

	return cardDataModel{
		Username:    c.Player.Username,
		Name:        c.Player.Name,
		Title:       string(c.Player.Title),
		Avatar:      c.Player.Avatar,
		Followers:   c.Player.Followers,
		CountryCode: c.Player.CountryCode,
		JoinedAt:    c.Player.JoinedAt,
		FideRating:  c.Stats.FideRating,
		Bullet:      toTimeControlModel(c.Stats.Bullet),
		Blitz:       toTimeControlModel(c.Stats.Blitz),
		Rapid:       toTimeControlModel(c.Stats.Rapid),
		Daily:       toTimeControlModel(c.Stats.Daily),
		OVR:         c.OVR,
		PlayStyle:   string(c.PlayStyle),
		Badges:      badges,
		TopOpenings: openings,
	}
}

func toTimeControlModel(s domain.TimeControlStats) timeControlStatsModel {
	return timeControlStatsModel{Rating: s.Rating, Highest: s.Highest, Wins: s.Wins, Losses: s.Losses, Draws: s.Draws}
}

func fromCardDataModel(m cardDataModel, cardType, tier string, computedAt, expiresAt time.Time) domain.Card {
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
			FideRating: m.FideRating,
			Bullet:     fromTimeControlModel(m.Bullet, domain.TimeControlBullet),
			Blitz:      fromTimeControlModel(m.Blitz, domain.TimeControlBlitz),
			Rapid:      fromTimeControlModel(m.Rapid, domain.TimeControlRapid),
			Daily:      fromTimeControlModel(m.Daily, domain.TimeControlDaily),
		},
		CardType:    domain.CardType(cardType),
		Tier:        domain.CardTier(tier),
		OVR:         m.OVR,
		PlayStyle:   domain.PlayStyle(m.PlayStyle),
		Badges:      badges,
		TopOpenings: openings,
		ComputedAt:  computedAt,
		ExpiresAt:   expiresAt,
	}
}

func fromTimeControlModel(m timeControlStatsModel, tc domain.TimeControl) domain.TimeControlStats {
	return domain.TimeControlStats{TimeControl: tc, Rating: m.Rating, Highest: m.Highest, Wins: m.Wins, Losses: m.Losses, Draws: m.Draws}
}
