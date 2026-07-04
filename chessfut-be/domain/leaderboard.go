package domain

type LeaderboardEntry struct {
	Username string
	Rank     int
	Rating   int
}

type Leaderboards struct {
	Bullet []LeaderboardEntry
	Blitz  []LeaderboardEntry
	Rapid  []LeaderboardEntry
	Daily  []LeaderboardEntry
}
