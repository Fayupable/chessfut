package service

const (
	chesscomMin = 100
	chesscomMax = 3500
)

func normalize(value int) float64 {
	if value <= chesscomMin {
		return 0
	}
	if value >= chesscomMax {
		return 100
	}
	return float64(value-chesscomMin) / float64(chesscomMax-chesscomMin) * 100
}
