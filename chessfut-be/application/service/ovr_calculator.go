package service

const (
	chesscomMin = 100
	chesscomMax = 3500
)

func normalize(value, min, max int) float64 {
	if value <= min {
		return 0
	}
	if value >= max {
		return 100
	}
	return float64(value-min) / float64(max-min) * 100
}
