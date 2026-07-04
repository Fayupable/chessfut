package http

import (
	"net/http"

	"github.com/fayupable/chessfut-be/application/port/input"
)

type AdminHandler struct {
	refreshCard       input.RefreshCardUseCase
	syncTitledPlayers input.SyncTitledPlayersUseCase
	refreshStaleCards input.RefreshStaleCardsUseCase
}

func NewAdminHandler(
	refreshCard input.RefreshCardUseCase,
	syncTitledPlayers input.SyncTitledPlayersUseCase,
	refreshStaleCards input.RefreshStaleCardsUseCase,
) *AdminHandler {
	return &AdminHandler{
		refreshCard:       refreshCard,
		syncTitledPlayers: syncTitledPlayers,
		refreshStaleCards: refreshStaleCards,
	}
}

func (h *AdminHandler) RefreshCard(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	card, err := h.refreshCard.Execute(r.Context(), username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toCardResponse(card))
}

func (h *AdminHandler) SyncTitledPlayers(w http.ResponseWriter, r *http.Request) {
	count, err := h.syncTitledPlayers.Execute(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"synced": count})
}

func (h *AdminHandler) RefreshStaleCards(w http.ResponseWriter, r *http.Request) {
	count, err := h.refreshStaleCards.Execute(r.Context(), 50)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"refreshed": count})
}
