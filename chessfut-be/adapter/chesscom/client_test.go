package chesscom

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fayupable/chessfut-be/domain"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/player/hikaru", r.URL.Path)
		assert.NotEmpty(t, r.Header.Get("User-Agent"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"username": "hikaru",
			"name": "Hikaru Nakamura",
			"title": "GM",
			"avatar": "https://example.com/avatar.png",
			"followers": 100,
			"country": "https://api.chess.com/pub/country/US",
			"joined": 1389043258
		}`))
	}))
	defer server.Close()

	client := NewClient("chessfut-test/1.0", server.URL)
	player, err := client.GetProfile(context.Background(), "hikaru")

	assert.NoError(t, err)
	assert.Equal(t, "hikaru", player.Username)
	assert.Equal(t, "GM", string(player.Title))
	assert.Equal(t, "US", player.CountryCode)
}

func TestClient_GetStats(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/player/hikaru/stats", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"chess_blitz": {"last": {"rating": 3414}, "best": {"rating": 3465}, "record": {"win": 100, "loss": 20, "draw": 10}},
			"fide": 2814
		}`))
	}))
	defer server.Close()

	client := NewClient("chessfut-test/1.0", server.URL)
	stats, err := client.GetStats(context.Background(), "hikaru")

	assert.NoError(t, err)
	assert.Equal(t, 2814, stats.FideRating)
	assert.Equal(t, 3414, stats.Blitz.Rating)
	assert.Equal(t, 100, stats.Blitz.Wins)
}

func TestClient_GetTitledUsernames(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/titled/GM", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"players": ["hikaru", "magnuscarlsen"]}`))
	}))
	defer server.Close()

	client := NewClient("chessfut-test/1.0", server.URL)
	usernames, err := client.GetTitledUsernames(context.Background(), "GM")

	assert.NoError(t, err)
	assert.Equal(t, []string{"hikaru", "magnuscarlsen"}, usernames)
}

func TestClient_DoRequest_HandlesNonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient("chessfut-test/1.0", server.URL)
	_, err := client.GetProfile(context.Background(), "nonexistent")

	assert.Error(t, err)
}

func TestClient_GetGames(t *testing.T) {
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/player/hikaru/games/archives":
			_, _ = w.Write([]byte(`{"archives": ["` + server.URL + `/player/hikaru/games/2024/01", "` + server.URL + `/player/hikaru/games/2020/01"]}`))
		case "/player/hikaru/games/2024/01":
			_, _ = w.Write([]byte(`{"games": [{
				"pgn": "[ECO \"B23\"]\n\n1. e4 c5 1-0",
				"time_control": "180",
				"end_time": 1704133706,
				"rated": true,
				"time_class": "blitz",
				"white": {"rating": 3236, "result": "win", "username": "hikaru"},
				"black": {"rating": 2867, "result": "resigned", "username": "someone"},
				"eco": "https://www.chess.com/openings/Sicilian-Defense"
			}]}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := NewClient("chessfut-test/1.0", server.URL)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	games, err := client.GetGames(context.Background(), "hikaru", from, to)

	assert.NoError(t, err)
	assert.Len(t, games, 1)
	assert.Equal(t, "B23", games[0].Opening.ECO)
	assert.Equal(t, domain.ColorWhite, games[0].Color)
	assert.Equal(t, domain.ResultWin, games[0].Result)
}
