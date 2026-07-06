package redis

import (
	"time"
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

type attributesModel struct {
	Pac int `json:"pac"`
	Sho int `json:"sho"`
	Pas int `json:"pas"`
	Dri int `json:"dri"`
	Def int `json:"def"`
	Phy int `json:"phy"`
}

type workRateModel struct {
	Attack  string `json:"attack"`
	Defense string `json:"defense"`
}

type cardModel struct {
	Username           string                `json:"username"`
	Name               string                `json:"name"`
	Title              string                `json:"title"`
	Avatar             string                `json:"avatar"`
	Followers          int                   `json:"followers"`
	CountryCode        string                `json:"country_code"`
	JoinedAt           time.Time             `json:"joined_at"`
	FideRating         int                   `json:"fide_rating"`
	TacticsRating      int                   `json:"tactics_rating"`
	PuzzleRushAccuracy float64               `json:"puzzle_rush_accuracy"`
	Bullet             timeControlStatsModel `json:"bullet"`
	Blitz              timeControlStatsModel `json:"blitz"`
	Rapid              timeControlStatsModel `json:"rapid"`
	Daily              timeControlStatsModel `json:"daily"`
	CardType           string                `json:"card_type"`
	Tier               string                `json:"tier"`
	OVR                int                   `json:"ovr"`
	PlayStyle          string                `json:"play_style"`
	Position           string                `json:"position"`
	Attributes         attributesModel       `json:"attributes"`
	WorkRate           workRateModel         `json:"work_rate"`
	Badges             []string              `json:"badges"`
	TopOpenings        []openingStatModel    `json:"top_openings"`
	GamesSnapshot      int                   `json:"games_snapshot"`
	ComputedAt         time.Time             `json:"computed_at"`
	ExpiresAt          time.Time             `json:"expires_at"`
}
