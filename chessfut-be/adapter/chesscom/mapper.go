package chesscom

import (
	"regexp"
	"strings"
	"time"

	"github.com/fayupable/chessfut-be/domain"
)

var (
	countryURLRe = regexp.MustCompile(`/country/([A-Z]{2,3})$`)
	ecoTagRe     = regexp.MustCompile(`\[ECO "([^"]+)"\]`)
	pgnTagRe     = regexp.MustCompile(`(?m)^\[.*\]\n?`)
	clockRe      = regexp.MustCompile(`\{[^}]*\}`)
	resultTokens = map[string]bool{"1-0": true, "0-1": true, "1/2-1/2": true, "*": true}
	moveNumberRe = regexp.MustCompile(`^\d+\.+$`)
)

func mapProfile(r profileResponse) domain.Player {
	return domain.Player{
		Username:    r.Username,
		Name:        r.Name,
		Title:       domain.Title(r.Title),
		Avatar:      r.Avatar,
		Followers:   r.Followers,
		JoinedAt:    time.Unix(r.Joined, 0),
		CountryCode: extractCountryCode(r.Country),
		URL:         r.URL,
		LastOnline:  time.Unix(r.LastOnline, 0),
	}
}

func extractCountryCode(countryURL string) string {
	match := countryURLRe.FindStringSubmatch(countryURL)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func mapStats(r statsResponse) domain.PlayerStats {
	return domain.PlayerStats{
		FideRating:         r.Fide,
		TacticsRating:      r.Tactics.Highest.Rating,
		PuzzleRushAccuracy: puzzleRushAccuracy(r.PuzzleRush),
		Bullet:             mapTimeControlStats(r.ChessBullet, domain.TimeControlBullet),
		Blitz:              mapTimeControlStats(r.ChessBlitz, domain.TimeControlBlitz),
		Rapid:              mapTimeControlStats(r.ChessRapid, domain.TimeControlRapid),
		Daily:              mapTimeControlStats(r.ChessDaily, domain.TimeControlDaily),
	}
}

func puzzleRushAccuracy(r puzzleRushResponse) float64 {
	if r.Best.TotalAttempts == 0 {
		return 0
	}
	return float64(r.Best.Score) / float64(r.Best.TotalAttempts) * 100
}

func mapTimeControlStats(r statsRecord, tc domain.TimeControl) domain.TimeControlStats {
	return domain.TimeControlStats{
		TimeControl: tc,
		Rating:      r.Last.Rating,
		Highest:     r.Best.Rating,
		Wins:        r.Record.Win,
		Losses:      r.Record.Loss,
		Draws:       r.Record.Draw,
		RD:          r.Last.RD,
	}
}

func mapGame(r gameResponse, forUsername string) domain.Game {
	color := domain.ColorWhite
	self, opponent := r.White, r.Black
	if !strings.EqualFold(r.White.Username, forUsername) {
		color = domain.ColorBlack
		self, opponent = r.Black, r.White
	}

	return domain.Game{
		URL:         "",
		PlayedAt:    time.Unix(r.EndTime, 0),
		TimeControl: mapTimeClass(r.TimeClass),
		Rated:       r.Rated,
		Color:       color,
		Result:      mapResult(self.Result),
		Opening:     mapOpening(r.ECO, r.PGN),
		MovesCount:  countMoves(r.PGN),
		PlayerElo:   self.Rating,
		OpponentElo: opponent.Rating,
	}
}

func mapTimeClass(tc string) domain.TimeControl {
	switch tc {
	case "bullet":
		return domain.TimeControlBullet
	case "blitz":
		return domain.TimeControlBlitz
	case "rapid":
		return domain.TimeControlRapid
	default:
		return domain.TimeControlDaily
	}
}

func mapResult(rawResult string) domain.GameResult {
	switch rawResult {
	case "win":
		return domain.ResultWin
	case "agreed", "repetition", "stalemate", "insufficient", "50move", "timevsinsufficient":
		return domain.ResultDraw
	default:
		return domain.ResultLoss
	}
}

func mapOpening(ecoURL, pgn string) domain.Opening {
	name := extractOpeningName(ecoURL)
	code := extractECOCode(pgn)
	return domain.Opening{ECO: code, Name: name}
}

func extractECOCode(pgn string) string {
	match := ecoTagRe.FindStringSubmatch(pgn)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func extractOpeningName(ecoURL string) string {
	parts := strings.Split(ecoURL, "/")
	if len(parts) == 0 {
		return ""
	}
	slug := parts[len(parts)-1]
	tokens := strings.Split(slug, "-")

	nameTokens := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		if len(tok) > 0 && (tok[0] >= '0' && tok[0] <= '9') {
			break
		}
		nameTokens = append(nameTokens, tok)
	}
	return strings.Join(nameTokens, " ")
}

func countMoves(pgn string) int {
	movetext := pgnTagRe.ReplaceAllString(pgn, "")
	movetext = clockRe.ReplaceAllString(movetext, "")

	tokens := strings.Fields(movetext)
	plies := 0
	for _, tok := range tokens {
		if resultTokens[tok] || moveNumberRe.MatchString(tok) {
			continue
		}
		plies++
	}
	return (plies + 1) / 2
}

func mapTitledUsernames(r titledResponse) []string {
	return r.Players
}
