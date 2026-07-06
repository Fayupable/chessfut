package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/fayupable/chessfut-be/application/port/input"
)

const (
	defaultLeaderboardLimit = 20
	maxLeaderboardLimit     = 50
	defaultSearchLimit      = 10
	maxSearchLimit          = 50
)

type CardHandler struct {
	getFastCard     input.GetFastCardUseCase
	getDetailedCard input.GetDetailedCardUseCase
	getLeaderboard  input.GetLeaderboardUseCase
	getStats        input.GetStatsUseCase
	searchPlayers   input.SearchPlayersUseCase
}

func NewCardHandler(
	getFastCard input.GetFastCardUseCase,
	getDetailedCard input.GetDetailedCardUseCase,
	getLeaderboard input.GetLeaderboardUseCase,
	getStats input.GetStatsUseCase,
	searchPlayers input.SearchPlayersUseCase,
) *CardHandler {
	return &CardHandler{
		getFastCard:     getFastCard,
		getDetailedCard: getDetailedCard,
		getLeaderboard:  getLeaderboard,
		getStats:        getStats,
		searchPlayers:   searchPlayers,
	}
}

func (h *CardHandler) GetFastCard(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	card, err := h.getFastCard.Execute(r.Context(), username)
	if err != nil {
		writeError(w, http.StatusNotFound, err, "player not found")
		return
	}

	writeJSON(w, http.StatusOK, toCardResponse(card))
}

func (h *CardHandler) GetDetailedCard(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	card, err := h.getDetailedCard.Execute(r.Context(), username)
	if err != nil {
		writeError(w, http.StatusNotFound, err, "player not found")
		return
	}

	writeJSON(w, http.StatusOK, toCardResponse(card))
}

func (h *CardHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	limit := parseBoundedLimit(r, defaultLeaderboardLimit, maxLeaderboardLimit)
	offset := parseOffset(r)

	cards, err := h.getLeaderboard.Execute(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err, "failed to load leaderboard")
		return
	}

	responses := make([]cardResponse, 0, len(cards))
	for _, card := range cards {
		responses = append(responses, toCardResponse(card))
	}

	writeJSON(w, http.StatusOK, map[string]any{"leaderboard": responses})
}

func (h *CardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	total, err := h.getStats.Execute(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err, "failed to load stats")
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"total_cards": total})
}

func (h *CardHandler) SearchPlayers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeJSON(w, http.StatusOK, map[string]any{"players": []string{}})
		return
	}

	limit := parseBoundedLimit(r, defaultSearchLimit, maxSearchLimit)
	offset := parseOffset(r)

	cards, err := h.searchPlayers.Execute(r.Context(), query, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err, "search failed")
		return
	}

	usernames := make([]string, 0, len(cards))
	for _, card := range cards {
		usernames = append(usernames, card.Player.Username)
	}

	writeJSON(w, http.StatusOK, map[string]any{"players": usernames})
}

func parseBoundedLimit(r *http.Request, defaultLimit, maxLimit int) int {
	limit := defaultLimit
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	return limit
}

func parseOffset(r *http.Request) int {
	offset := 0
	if raw := r.URL.Query().Get("offset"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed >= 0 {
			offset = parsed
		}
	}
	return offset
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// writeError logs the real, potentially sensitive error server-side and
// returns only a fixed public-facing message to the client.
func writeError(w http.ResponseWriter, status int, err error, publicMessage string) {
	slog.Error("request_failed", "status", status, "error", err.Error())
	writeJSON(w, status, map[string]string{"error": publicMessage})
}
