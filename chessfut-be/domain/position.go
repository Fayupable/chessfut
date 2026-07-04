package domain

type Position string

const (
	PositionForward    Position = "ST"
	PositionMidfielder Position = "CM"
	PositionDefender   Position = "CB"
	PositionAllRounder Position = "CAM"
)

func MapToPosition(style PlayStyle) Position {
	switch style {
	case PlayStyleAggressive:
		return PositionForward
	case PlayStylePositional:
		return PositionMidfielder
	case PlayStyleDefensive:
		return PositionDefender
	default:
		return PositionAllRounder
	}
}
