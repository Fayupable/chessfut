package service

const (
	chesscomMin = 100
	chesscomMax = 3500
)

func normalize(value, max int) float64 {
	if value <= chesscomMin {
		return 0
	}
	if value >= max {
		return 100
	}
	return float64(value-chesscomMin) / float64(max-chesscomMin) * 100
}
