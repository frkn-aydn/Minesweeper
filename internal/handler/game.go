package handlers

import (
	"mine-game/internal/services"
	"net/http"
)

type GameHandler struct {
	*Helper
	gameService *services.GameService
}

// NewGameHandler creates a new game handler.
func NewGameHandler(gameService *services.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

// CreateGame creates a new game.
func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	mines, err := h.StringToInt(r.URL.Query().Get("mines"))
	if err != nil {
		h.WriteErrorResponse(w, http.StatusBadRequest, "mines is required", err.Error())
		return
	}

	game, err := h.gameService.CreateNewGame(ctx, mines)
	if err != nil {
		h.WriteErrorResponse(w, http.StatusInternalServerError, "failed to create game", err.Error())
		return
	}

	h.WriteSuccessResponse(w, http.StatusCreated, game)
}

// MakeMove makes a move in the game.
func (h *GameHandler) MakeMove(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	gameID, err := h.StringToPrimitiveObjectID(r.URL.Query().Get("game_id"))
	if err != nil || gameID.IsZero() {
		h.WriteErrorResponse(w, http.StatusBadRequest, "game_id is required", err.Error())
		return
	}

	x, err := h.StringToInt(r.URL.Query().Get("x"))
	if err != nil {
		h.WriteErrorResponse(w, http.StatusBadRequest, "x is required", err.Error())
		return
	}

	y, err := h.StringToInt(r.URL.Query().Get("y"))
	if err != nil {
		h.WriteErrorResponse(w, http.StatusBadRequest, "y is required", err.Error())
		return
	}

	game, err := h.gameService.GetGame(ctx, gameID)
	if err != nil {
		h.WriteErrorResponse(w, http.StatusInternalServerError, "failed to get game:"+gameID.Hex(), err.Error())
		return
	}

	if game.GameOver {
		h.WriteErrorResponse(w, http.StatusOK, "game is over", nil)
		return
	}

	game, err = game.MakeMove(x, y)
	if err != nil {
		h.WriteErrorResponse(w, http.StatusOK, err.Error(), nil)
		return
	}

	err = h.gameService.UpdateGame(ctx, game)
	if err != nil {
		h.WriteErrorResponse(w, http.StatusInternalServerError, "failed to update game", err.Error())
		return
	}

	if game.GameOver {
		h.WriteSuccessResponse(w, http.StatusOK, map[string]interface{}{
			"board":     game.Board,
			"payout":    game.Payout,
			"game_over": game.GameOver,
			"id":        game.ID,
			"mines":     game.Mines,
			"moves":     game.Moves,
			"revealed":  game.Revealed,
		})
		return
	}

	h.WriteSuccessResponse(w, http.StatusCreated, game)
}
