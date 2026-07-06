package chesscom

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/sync/singleflight"
	"golang.org/x/time/rate"

	"github.com/fayupable/chessfut-be/application/port/output"
	"github.com/fayupable/chessfut-be/domain"
)

type Client struct {
	httpClient *http.Client
	limiter    *rate.Limiter
	group      singleflight.Group
	userAgent  string
	baseURL    string
}

func (c *Client) GetGames(ctx context.Context, username string, from, to time.Time) ([]domain.Game, error) {
	var archives archivesResponse
	archivesURL := fmt.Sprintf("%s/player/%s/games/archives", c.baseURL, url.PathEscape(username))
	if err := c.doRequest(ctx, archivesURL, &archives); err != nil {
		return nil, err
	}

	var games []domain.Game
	for _, archiveURL := range archives.Archives {
		archiveDate, ok := parseArchiveDate(archiveURL)
		if !ok || archiveDate.Before(monthFloor(from)) || archiveDate.After(to) {
			continue
		}

		var resp gamesResponse
		if err := c.doRequest(ctx, archiveURL, &resp); err != nil {
			return nil, err
		}

		for _, g := range resp.Games {
			game := mapGame(g, username)
			if game.PlayedAt.Before(from) || game.PlayedAt.After(to) {
				continue
			}
			games = append(games, game)
		}
	}

	return games, nil
}

func NewClient(userAgent, baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		limiter:    rate.NewLimiter(rate.Limit(1), 2),
		userAgent:  userAgent,
		baseURL:    baseURL,
	}
}

var _ output.ChessComClientPort = (*Client)(nil)

func (c *Client) doRequest(ctx context.Context, url string, target any) error {
	start := time.Now()

	if err := c.limiter.Wait(ctx); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error("chesscom_request", "url", url, "error", err.Error(), "duration_ms", time.Since(start).Milliseconds())
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	slog.Info("chesscom_request",
		"url", url,
		"status", resp.StatusCode,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	if resp.StatusCode == http.StatusTooManyRequests {
		c.backOff()
		return fmt.Errorf("chesscom: rate limited (429) for %s", url)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("chesscom: unexpected status %d for %s", resp.StatusCode, url)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func (c *Client) backOff() {
	c.limiter.SetLimit(rate.Limit(0.2))
	time.AfterFunc(2*time.Minute, func() {
		c.limiter.SetLimit(rate.Limit(1))
	})
}

func (c *Client) fetchDeduped(ctx context.Context, key, url string, target any) error {
	_, err, _ := c.group.Do(key, func() (any, error) {
		return nil, c.doRequest(ctx, url, target)
	})
	return err
}

func (c *Client) GetProfile(ctx context.Context, username string) (domain.Player, error) {
	var resp profileResponse
	reqURL := fmt.Sprintf("%s/player/%s", c.baseURL, url.PathEscape(username))
	if err := c.fetchDeduped(ctx, "profile:"+username, reqURL, &resp); err != nil {
		return domain.Player{}, err
	}
	return mapProfile(resp), nil
}

func (c *Client) GetStats(ctx context.Context, username string) (domain.PlayerStats, error) {
	var resp statsResponse
	reqURL := fmt.Sprintf("%s/player/%s/stats", c.baseURL, url.PathEscape(username))
	if err := c.fetchDeduped(ctx, "stats:"+username, reqURL, &resp); err != nil {
		return domain.PlayerStats{}, err
	}
	return mapStats(resp), nil
}

func (c *Client) GetTitledUsernames(ctx context.Context, title domain.Title) ([]string, error) {
	var resp titledResponse
	url := fmt.Sprintf("%s/titled/%s", c.baseURL, title)
	if err := c.doRequest(ctx, url, &resp); err != nil {
		return nil, err
	}
	return mapTitledUsernames(resp), nil
}
