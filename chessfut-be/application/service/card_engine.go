package service

import (
	"math"

	"github.com/fayupable/chessfut-be/domain"
)

const (
	idxPac = iota
	idxSho
	idxPas
	idxDri
	idxDef
	idxPhy
	statCount
)

const (
	elasticAlpha = 0.4
	statFloor    = 10.0
	statCeiling  = 99.0

	fideX0    = 2500.0
	fideK     = 0.0084
	fideBase  = 50.0
	fideRange = 49.0

	chesscomX0    = 2600.0
	chesscomK     = 0.0045
	chesscomBase  = 30.0
	chesscomRange = 65.0

	fideWeight = 0.8

	longGameMoves = 80.0

	maxAnchorDeviation = 0.45

	inactivityFactor = 0.65
)

const maxPlausibleFideRating = 2900

func clampFloat(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func mean(v []float64) float64 {
	sum := 0.0
	for _, x := range v {
		sum += x
	}
	return sum / float64(len(v))
}

// BuildCardScoring anchors the whole scoring pipeline to a player's *verified*
// strength: FIDE rating when available (the authoritative, tightly-banded
// real-world measure), falling back to chess.com's own rating scale otherwise.
// Every attribute is elastically shaped around that anchor, so two GMs with
// genuinely different real-world strength (e.g. FIDE 2500 vs 2800) land at
// meaningfully different OVRs instead of both saturating at the cap.
func BuildCardScoring(
	player domain.Player,
	stats domain.PlayerStats,
	topOpenings []domain.OpeningStat,
	games []domain.Game,
	gamesSnapshot int,
) (domain.Attributes, domain.Position, int, domain.WorkRate) {
	raw := rawAttributeValues(player.Title, stats, topOpenings, gamesSnapshot)

	center := anchorScore(player.Title, stats)
	attrs := elasticShape(raw, center)

	drawRate := blitzDrawRate(stats.Blitz)
	aggression := aggressionSignal(games)

	position, family := positionFromVectors(attrs, drawRate, aggression, avgGameLength(games))
	weightedScore := weightedOVRScore(attrs, family)
	// Position weighting can zero out an attribute entirely (e.g. forwards
	// ignore DEF/PHY), which lets a weak stat inflate the *other* stats via
	// elasticShape's mean-relative shaping without ever being penalized —
	// producing a final OVR that contradicts the anchor's verified ranking.
	// Bounding the weighted score to a small band around the anchor keeps
	// position flavor while guaranteeing the anchor's ranking always holds.
	boundedScore := clampFloat(weightedScore, center-maxAnchorDeviation, center+maxAnchorDeviation)

	nudge := titleNudge(player.Title)
	ovr := int(clampFloat(math.Round(boundedScore+nudge*(99-boundedScore)/99), 1, 99))

	workRate := CalculateWorkRate(attrs)

	return attrs, position, ovr, workRate
}

func rawAttributeValues(title domain.Title, stats domain.PlayerStats, topOpenings []domain.OpeningStat, gamesSnapshot int) [statCount]float64 {
	var raw [statCount]float64

	speed := speedSignal(stats)
	raw[idxPac] = clampFloat(10+(speed-100)/30, statFloor, statCeiling)

	raw[idxSho] = clampFloat(stats.Blitz.WinRate()*1.5+10, statFloor, statCeiling)

	raw[idxPas] = clampFloat(10+(passSignal(stats)-100)/30, statFloor, statCeiling)

	raw[idxDri] = technicalScore(stats, topOpenings)

	total := stats.Blitz.TotalGames()
	nonLoss := 0.0
	if total > 0 {
		nonLoss = float64(stats.Blitz.Wins+stats.Blitz.Draws) / float64(total) * 130
	}
	raw[idxDef] = clampFloat(nonLoss, statFloor, statCeiling)

	volumePhy := math.Log(float64(gamesSnapshot)+1) / math.Log(100000) * 99
	raw[idxPhy] = clampFloat(math.Max(volumePhy, titlePhysicalFloor(title)), statFloor, statCeiling)

	return raw
}

// speedSignal falls back to blitz rating alone when bullet was never played,
// so an absent format is treated as missing data rather than a real-zero
// rating that would otherwise crater the PAC score.
func speedSignal(stats domain.PlayerStats) float64 {
	if stats.Bullet.Rating == 0 {
		return float64(stats.Blitz.Rating)
	}
	return 0.6*float64(stats.Bullet.Rating) + 0.4*float64(stats.Blitz.Rating)
}

// passSignal falls back to the faster time controls when rapid was never
// played, for the same reason as speedSignal.
func passSignal(stats domain.PlayerStats) float64 {
	if stats.Rapid.Rating > 0 {
		return float64(stats.Rapid.Rating)
	}
	return math.Max(float64(stats.Bullet.Rating), float64(stats.Blitz.Rating))
}

// titlePhysicalFloor keeps a newly-created titled account from being
// penalized on PHY (a pure games-played signal) just for not having racked
// up a large game count yet.
func titlePhysicalFloor(title domain.Title) float64 {
	switch title {
	case domain.TitleGM, domain.TitleWGM:
		return 70
	case domain.TitleIM, domain.TitleWIM:
		return 65
	case domain.TitleFM, domain.TitleWFM:
		return 60
	case domain.TitleCM, domain.TitleWCM:
		return 55
	default:
		return 0
	}
}

func technicalScore(stats domain.PlayerStats, topOpenings []domain.OpeningStat) float64 {
	components := []float64{openingDiversity(topOpenings)}

	if stats.TacticsRating > 0 {
		components = append(components, normalize(stats.TacticsRating))
	}
	if stats.PuzzleRushAccuracy > 0 {
		components = append(components, clampFloat(stats.PuzzleRushAccuracy, 0, 99))
	}

	return clampFloat(mean(components), statFloor, statCeiling)
}

func peakRating(stats domain.PlayerStats) float64 {
	peak := math.Max(float64(stats.Bullet.Rating), float64(stats.Blitz.Rating))
	return math.Max(peak, float64(stats.Rapid.Rating))
}

// FideRatingInfo reports both the FIDE rating actually used for scoring and
// where it came from, so callers (the API response, the frontend) can
// explain the number instead of presenting it as an unquestionable fact.
//   - "verified": chess.com returned a real, plausible fide_rating
//   - "title_default": fide_rating was missing or implausible, so the
//     title's minimum norm rating was assumed instead
//   - "none": no title and no fide_rating — chess.com rating alone drives
//     the score
func FideRatingInfo(title domain.Title, fideRating int) (effective int, source string) {
	if fideRating > 0 && fideRating <= maxPlausibleFideRating {
		return fideRating, "verified"
	}

	switch title {
	case domain.TitleGM, domain.TitleWGM:
		return 2500, "title_default"
	case domain.TitleIM, domain.TitleWIM:
		return 2400, "title_default"
	case domain.TitleFM, domain.TitleWFM:
		return 2300, "title_default"
	case domain.TitleCM, domain.TitleWCM:
		return 2200, "title_default"
	default:
		return 0, "none"
	}
}

// effectiveFideRating falls back to the title's minimum norm rating when
// chess.com doesn't have the player's FIDE rating linked, so a verified title
// still counts for something instead of being treated as fully unrated. A
// fide_rating above the real-world record (Carlsen's 2882) is untrustworthy
// self-reported data, not a genuine rating — it's ignored entirely rather
// than trusted-but-capped.
func effectiveFideRating(title domain.Title, fideRating int) int {
	effective, _ := FideRatingInfo(title, fideRating)
	return effective
}

// anchorScore is the single "how strong is this player, really" number that
// everything else (attribute shaping, OVR) is built around. FIDE dominates
// when present (or assumed from title) since it's independently verified and
// tightly banded; chess.com activity contributes a smaller nudge on top.
// A player with zero games in every chess.com format has no platform-verified
// performance to blend in — their FIDE-derived score is discounted by
// inactivityFactor instead of granted in full, so an inactive titled account
// never outranks someone who's actually proven their strength on chess.com.
// Without FIDE or a title, chess.com rating alone drives it — no hard
// ceiling, a genuinely strong untitled player can still land high.
func anchorScore(title domain.Title, stats domain.PlayerStats) float64 {
	fide := effectiveFideRating(title, stats.FideRating)
	peak := peakRating(stats)

	if peak <= 0 {
		if fide <= 0 {
			return statFloor
		}
		fideScore := fideBase + fideRange/(1+math.Exp(-fideK*(float64(fide)-fideX0)))
		return math.Max(fideScore*inactivityFactor, statFloor)
	}

	chesscomScore := chesscomBase + chesscomRange/(1+math.Exp(-chesscomK*(peak-chesscomX0)))
	if fide <= 0 {
		return chesscomScore
	}

	fideScore := fideBase + fideRange/(1+math.Exp(-fideK*(float64(fide)-fideX0)))
	return fideWeight*fideScore + (1-fideWeight)*chesscomScore
}

func elasticShape(raw [statCount]float64, center float64) domain.Attributes {
	m := mean(raw[:])

	shape := func(x float64) int {
		return int(math.Round(clampFloat(center+elasticAlpha*(x-m), statFloor, statCeiling)))
	}

	return domain.Attributes{
		Pac: shape(raw[idxPac]),
		Sho: shape(raw[idxSho]),
		Pas: shape(raw[idxPas]),
		Dri: shape(raw[idxDri]),
		Def: shape(raw[idxDef]),
		Phy: shape(raw[idxPhy]),
	}
}

func blitzDrawRate(tc domain.TimeControlStats) float64 {
	total := tc.TotalGames()
	if total == 0 {
		return 0
	}
	return float64(tc.Draws) / float64(total)
}

func avgGameLength(games []domain.Game) float64 {
	if len(games) == 0 {
		return longGameMoves / 2
	}
	total := 0
	for _, g := range games {
		total += g.MovesCount
	}
	return float64(total) / float64(len(games))
}

func aggressionSignal(games []domain.Game) float64 {
	if len(games) == 0 {
		return 0.5
	}

	decisive := 0
	for _, g := range games {
		if g.Result != domain.ResultDraw {
			decisive++
		}
	}
	decisiveRate := float64(decisive) / float64(len(games))

	shortGame := clampFloat(1-avgGameLength(games)/60, 0, 1)

	return 0.5*shortGame + 0.5*decisiveRate
}

func positionFromVectors(attrs domain.Attributes, drawRate, aggression, avgMoves float64) (domain.Position, string) {
	attack := 0.5*float64(attrs.Sho) + 0.3*float64(attrs.Pac) + 0.2*(aggression*99)
	playmaker := 0.4*float64(attrs.Pas) + 0.4*float64(attrs.Dri) + 0.2*((1-drawRate)*99)
	anchor := 0.6*float64(attrs.Def) + 0.2*(drawRate*99) + 0.2*((1-clampFloat(avgMoves/longGameMoves, 0, 1))*99)

	switch {
	case attack >= playmaker && attack >= anchor:
		return domain.PositionForward, "forward"
	case playmaker >= anchor:
		if attrs.Dri >= attrs.Pas {
			return domain.PositionAllRounder, "playmaker"
		}
		return domain.PositionMidfielder, "playmaker"
	default:
		return domain.PositionDefender, "anchor"
	}
}

func weightedOVRScore(attrs domain.Attributes, family string) float64 {
	weights := map[string][statCount]float64{
		"forward":   {0.30, 0.40, 0.10, 0.20, 0.00, 0.00},
		"playmaker": {0.20, 0.00, 0.30, 0.30, 0.20, 0.00},
		"anchor":    {0.00, 0.00, 0.20, 0.00, 0.50, 0.30},
	}[family]

	vals := [statCount]float64{
		float64(attrs.Pac), float64(attrs.Sho), float64(attrs.Pas),
		float64(attrs.Dri), float64(attrs.Def), float64(attrs.Phy),
	}

	sum := 0.0
	for i, w := range weights {
		sum += vals[i] * w
	}

	return sum
}

func titleNudge(title domain.Title) float64 {
	switch title {
	case domain.TitleGM, domain.TitleWGM:
		return 3
	case domain.TitleIM, domain.TitleWIM:
		return 2
	case domain.TitleFM, domain.TitleWFM:
		return 1
	case domain.TitleCM, domain.TitleWCM:
		return 0.5
	default:
		return 0
	}
}
