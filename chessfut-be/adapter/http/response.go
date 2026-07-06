package http

import (
	"time"
)

type timeControlResponse struct {
	Rating  int     `json:"rating"`
	Highest int     `json:"highest"`
	Wins    int     `json:"wins"`
	Losses  int     `json:"losses"`
	Draws   int     `json:"draws"`
	WinRate float64 `json:"win_rate"`
}

type openingResponse struct {
	ECO     string  `json:"eco"`
	Name    string  `json:"name"`
	Count   int     `json:"count"`
	WinRate float64 `json:"win_rate"`
}

type attributesResponse struct {
	Pac int `json:"pac"`
	Sho int `json:"sho"`
	Pas int `json:"pas"`
	Dri int `json:"dri"`
	Def int `json:"def"`
	Phy int `json:"phy"`
}

type workRateResponse struct {
	Attack  string `json:"attack"`
	Defense string `json:"defense"`
}

type cardResponse struct {
	Username           string              `json:"username"`
	Name               string              `json:"name"`
	Title              string              `json:"title"`
	Avatar             string              `json:"avatar"`
	Followers          int                 `json:"followers"`
	CountryCode        string              `json:"country_code"`
	JoinedAt           time.Time           `json:"joined_date"`
	CardType           string              `json:"card_type"`
	Tier               string              `json:"tier"`
	OVR                int                 `json:"ovr"`
	PlayStyle          string              `json:"play_style,omitempty"`
	Position           string              `json:"position,omitempty"`
	TacticsRating      int                 `json:"tactics_rating,omitempty"`
	PuzzleRushAccuracy float64             `json:"puzzle_rush_accuracy,omitempty"`
	Attributes         attributesResponse  `json:"attributes"`
	WorkRate           workRateResponse    `json:"work_rate"`
	Badges             []string            `json:"badges"`
	Bullet             timeControlResponse `json:"bullet"`
	Blitz              timeControlResponse `json:"blitz"`
	Rapid              timeControlResponse `json:"rapid"`
	Daily              timeControlResponse `json:"daily"`
	TopOpenings        []openingResponse   `json:"top_openings,omitempty"`
	LastUpdated        time.Time           `json:"last_updated"`
}
