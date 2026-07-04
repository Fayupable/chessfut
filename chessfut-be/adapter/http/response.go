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

type cardResponse struct {
	Username    string              `json:"username"`
	Name        string              `json:"name"`
	Title       string              `json:"title"`
	Avatar      string              `json:"avatar"`
	Followers   int                 `json:"followers"`
	CountryCode string              `json:"country_code"`
	JoinedAt    time.Time           `json:"joined_date"`
	CardType    string              `json:"card_type"`
	Tier        string              `json:"tier"`
	OVR         int                 `json:"ovr"`
	PlayStyle   string              `json:"play_style,omitempty"`
	Position    string              `json:"position,omitempty"`
	Badges      []string            `json:"badges"`
	Bullet      timeControlResponse `json:"bullet"`
	Blitz       timeControlResponse `json:"blitz"`
	Rapid       timeControlResponse `json:"rapid"`
	Daily       timeControlResponse `json:"daily"`
	TopOpenings []openingResponse   `json:"top_openings,omitempty"`
	LastUpdated time.Time           `json:"last_updated"`
}
