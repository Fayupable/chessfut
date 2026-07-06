package domain

import "time"

type Title string

const (
	TitleGM   Title = "GM"
	TitleWGM  Title = "WGM"
	TitleIM   Title = "IM"
	TitleWIM  Title = "WIM"
	TitleFM   Title = "FM"
	TitleWFM  Title = "WFM"
	TitleCM   Title = "CM"
	TitleWCM  Title = "WCM"
	TitleNone Title = ""
)

type Player struct {
	Username    string
	Name        string
	Title       Title
	Avatar      string
	Followers   int
	JoinedAt    time.Time
	CountryCode string
}

func (p Player) HasFideTitle() bool {
	return p.Title != TitleNone
}

type TimeControl string

const (
	TimeControlBullet TimeControl = "bullet"
	TimeControlBlitz  TimeControl = "blitz"
	TimeControlRapid  TimeControl = "rapid"
	TimeControlDaily  TimeControl = "daily"
)

type TimeControlStats struct {
	TimeControl TimeControl
	Rating      int
	Highest     int
	Wins        int
	Losses      int
	Draws       int
}

func (s TimeControlStats) TotalGames() int {
	return s.Wins + s.Losses + s.Draws
}

func (s TimeControlStats) WinRate() float64 {
	total := s.TotalGames()
	if total == 0 {
		return 0
	}
	return float64(s.Wins) / float64(total) * 100
}

type PlayerStats struct {
	FideRating         int
	TacticsRating      int
	PuzzleRushAccuracy float64
	Bullet             TimeControlStats
	Blitz              TimeControlStats
	Rapid              TimeControlStats
	Daily              TimeControlStats
}
