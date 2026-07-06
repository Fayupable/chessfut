package domain

type Position string

const (
	PositionForward    Position = "ST"
	PositionMidfielder Position = "CM"
	PositionDefender   Position = "CB"
	PositionAllRounder Position = "CAM"
)
