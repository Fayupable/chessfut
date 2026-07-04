package domain

import "time"

type Color string

const (
	ColorWhite Color = "white"
	ColorBlack Color = "black"
)

type GameResult string

const (
	ResultWin  GameResult = "win"
	ResultLoss GameResult = "loss"
	ResultDraw GameResult = "draw"
)

type Opening struct {
	ECO  string
	Name string
}

type Game struct {
	URL         string
	PlayedAt    time.Time
	TimeControl TimeControl
	Rated       bool
	Color       Color
	Result      GameResult
	Opening     Opening
	MovesCount  int
	PlayerElo   int
	OpponentElo int
}

type OpeningStat struct {
	Opening Opening
	Count   int
	Wins    int
	Losses  int
	Draws   int
}

func (o OpeningStat) WinRate() float64 {
	if o.Count == 0 {
		return 0
	}
	return float64(o.Wins) / float64(o.Count) * 100
}
