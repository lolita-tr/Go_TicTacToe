package game

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	gameService Service
	logger      *zap.Logger
}

func NewGameHandler(gameService Service, logger *zap.Logger) *Handler {
	return &Handler{
		gameService: gameService,
		logger:      logger,
	}
}

func (gh *Handler) CreateBotGame(w http.ResponseWriter, r *http.Request) {
	userUUID, ok := r.Context().Value("user_id").(string)
	if !ok || userUUID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response, err := gh.gameService.CreateBotGame(userUUID)
	if err != nil {
		http.Error(w, "Failed to create game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (gh *Handler) CreateUserGame(w http.ResponseWriter, r *http.Request) {
	var req NewGameUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userUUID, ok := r.Context().Value("user_id").(string)
	if !ok || userUUID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response, err := gh.gameService.CreateUserGame(userUUID, req.OpponentLogin)
	if err != nil {
		http.Error(w, "Failed to create game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (gh *Handler) MakeMovePlayers(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("id")
	if gameID == "" {
		http.Error(w, "missing game ID", http.StatusBadRequest)
		return
	}

	var request MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	response, err := gh.gameService.MakeMovePlayers(request, gameID)
	if err != nil {
		http.Error(w, "Failed to make move", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (gh *Handler) MakeMoveBot(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("id")
	if gameID == "" {
		http.Error(w, "missing game ID", http.StatusBadRequest)
		return
	}

	var request MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	response, err := gh.gameService.MakeMoveBot(request, gameID)
	if err != nil {
		http.Error(w, "Failed to make move", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (gh *Handler) GetBoard(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("id")
	if gameID == "" || gameID == "undefined" {
		http.Error(w, "Missing game ID", http.StatusBadRequest)
		return
	}

	game, err := gh.gameService.getGame(gameID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(game.GetBoard())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (gh *Handler) ResetGame(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("id")
	if gameID == "" {
		http.Error(w, "Missing game ID", http.StatusBadRequest)
		return
	}

	game, err := gh.gameService.getGame(gameID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	game.Reset()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "reset",
	})
}

func (gh *Handler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UserInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.UserUUID == "" {
		http.Error(w, "UserUUID is required", http.StatusBadRequest)
		return
	}

	response, err := gh.gameService.GetUserInfo(context.Background(), req.UserUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (gh *Handler) GetCurrentGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	response := gh.gameService.GetCurrentGames()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (gh *Handler) GetGame(w http.ResponseWriter, r *http.Request) {
	var req GameInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	gameID := req.GameUUID
	response, err := gh.gameService.GetGameInfo(gameID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (gh *Handler) GetFinishedGames(w http.ResponseWriter, r *http.Request) {
	var req FinishedGamesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	response, err := gh.gameService.GetFinishedGames(context.Background(), req.UserUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (gh *Handler) GetTopPlayers(w http.ResponseWriter, r *http.Request) {
	var req TopPlayersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	response, err := gh.gameService.GetTopPlayers(context.Background(), req.Top)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
