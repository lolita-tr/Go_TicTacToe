package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"src/internal/domain"
	"sync"

	"go.uber.org/zap"
)

type GameHandler struct {
	logger *zap.Logger
	games  *sync.Map
}

func NewGameHandler(logger *zap.Logger) *GameHandler {
	return &GameHandler{
		logger: logger,
		games:  &sync.Map{},
	}
}

func (h *GameHandler) HandleNewGame(w http.ResponseWriter, r *http.Request) {
	id := uuid.New().String()
	newGame := domain.CreateGame()

	h.games.Store(id, newGame)

	json.NewEncoder(w).Encode(map[string]any{
		"game_id": id,
	})
}

func (h *GameHandler) HandleMove(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("id")
	if gameID == "" {
		http.Error(w, "Missing game ID", http.StatusBadRequest)
		return
	}

	value, ok := h.games.Load(gameID)
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}
	game := value.(*domain.Game)

	var req struct {
		Row int `json:"row"`
		Col int `json:"col"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	win, err := game.MakeMove(req.Row, req.Col)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if win {
		json.NewEncoder(w).Encode(map[string]any{
			"winner": game.CurrentPlayer(),
			"board":  game.GetBoard(),
		})
		return
	}

	if !game.Table.IsFull() {
		row, col := game.Table.FindBestMove()
		winRobot, errRobot := game.MakeMove(row, col)
		if errRobot != nil {
			http.Error(w, errRobot.Error(), http.StatusBadRequest)
			return
		}

		if winRobot {
			json.NewEncoder(w).Encode(map[string]any{
				"winner": game.CurrentPlayer(),
				"board":  game.GetBoard(),
			})
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]any{
		"winner": 0,
		"board":  game.GetBoard(),
	})
}

func (h *GameHandler) HandleBoard(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("id")
	if gameID == "" {
		http.Error(w, "Missing game ID", http.StatusBadRequest)
		return
	}

	value, ok := h.games.Load(gameID)
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}
	game := value.(*domain.Game)

	json.NewEncoder(w).Encode(game.GetBoard())
}

func (h *GameHandler) HandleReset(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("id")
	if gameID == "" {
		http.Error(w, "Missing game ID", http.StatusBadRequest)
		return
	}

	value, ok := h.games.Load(gameID)
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}
	game := value.(*domain.Game)

	game.Reset()

	json.NewEncoder(w).Encode(map[string]any{
		"status": "reset",
	})
}
