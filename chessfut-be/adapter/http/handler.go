package http

import (
	"encoding/json"
	"net/http"

	"github.com/fayupable/chessfut-be/application/port/input"
)

type CardHandler struct {
	getFastCard     input.GetFastCardUseCase
	getDetailedCard input.GetDetailedCardUseCase
	getLeaderboard  input.GetLeaderboardUseCase
}

func NewCardHandler(
	getFastCard input.GetFastCardUseCase,
	getDetailedCard input.GetDetailedCardUseCase,
	getLeaderboard input.GetLeaderboardUseCase,
) *CardHandler {
	return &CardHandler{getFastCard: getFastCard, getDetailedCard: getDetailedCard, getLeaderboard: getLeaderboard}
}

func (h *CardHandler) GetFastCard(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	card, err := h.getFastCard.Execute(r.Context(), username)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toCardResponse(card))
}

func (h *CardHandler) GetDetailedCard(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	card, err := h.getDetailedCard.Execute(r.Context(), username)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toCardResponse(card))
}

func (h *CardHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	limit := 50

	cards, err := h.getLeaderboard.Execute(r.Context(), limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]cardResponse, 0, len(cards))
	for _, card := range cards {
		responses = append(responses, toCardResponse(card))
	}

	writeJSON(w, http.StatusOK, map[string]any{"leaderboard": responses})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
