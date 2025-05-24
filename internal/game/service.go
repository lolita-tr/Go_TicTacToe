package game

import (
	"Go_TicTacToe/internal/domain"
	"Go_TicTacToe/internal/storage"
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"sync"
)

type ServiceImpl struct {
	userRepository storage.UserRepository
	gameRepository storage.GameRepository
	games          *sync.Map
}

func NewGameService(gameRepository storage.GameRepository, userRepository storage.UserRepository) Service {
	return &ServiceImpl{
		userRepository: userRepository,
		gameRepository: gameRepository,
		games:          &sync.Map{},
	}
}

func (gs *ServiceImpl) CreateBotGame(userID string) (*NewGameResponse, error) {
	gameID := uuid.New().String()

	newGame := domain.CreateGame()
	newGame.Mode = "bot"
	newGame.Player1 = userID

	gs.games.Store(gameID, newGame)

	if err := gs.gameRepository.SaveNewGame(context.Background(), gameID, newGame.Mode); err != nil {
		return nil, err
	}
	response := &NewGameResponse{
		gameID,
	}

	return response, nil
}

func (gs *ServiceImpl) CreateUserGame(userID string, oppLogin string) (*NewGameResponse, error) {
	gameID := uuid.New().String()

	oppID, err := gs.userRepository.GetUserUUIDbyLogin(context.Background(), oppLogin)

	newGame := domain.CreateGame()
	newGame.Mode = "pvp"
	newGame.Player1 = userID
	newGame.Player2 = oppID

	gs.games.Store(gameID, newGame)
	if err = gs.gameRepository.SaveNewGame(context.Background(), gameID, newGame.Mode); err != nil {
		return nil, errors.Wrap(err, "failed to save game")
	}

	response := &NewGameResponse{
		gameID,
	}

	return response, nil
}

func (gs *ServiceImpl) MakeMovePlayers(request MoveRequest, gameID string) (*MoveResponse, error) {
	game, err := gs.getGame(gameID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch game")
	}

	win, err := game.MakeMove(request.Row, request.Col)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make move")
	}

	if win {
		if game.GetPlayerUUID(game.CurrentPlayer()) == game.Player1 {
			gs.gameRepository.SaveResult(context.Background(), game.Player1, gameID, "win")
			gs.gameRepository.SaveResult(context.Background(), game.Player2, gameID, "lose")
		} else if game.GetPlayerUUID(game.CurrentPlayer()) == game.Player2 {
			gs.gameRepository.SaveResult(context.Background(), game.Player1, gameID, "lose")
			gs.gameRepository.SaveResult(context.Background(), game.Player2, gameID, "win")
		}

		response := &MoveResponse{
			Board:  game.GetBoard(),
			Result: game.GetPlayerUUID(game.CurrentPlayer()),
		}
		return response, nil
	}

	if game.Table.IsFull() {
		gs.gameRepository.SaveResult(context.Background(), game.Player1, gameID, "draw")
		gs.gameRepository.SaveResult(context.Background(), game.Player2, gameID, "draw")

		response := &MoveResponse{
			Board:  game.GetBoard(),
			Result: "draw",
		}
		return response, nil
	}

	response := &MoveResponse{
		Board:  game.GetBoard(),
		Result: "",
	}
	return response, nil
}

func (gs *ServiceImpl) MakeMoveBot(request MoveRequest, gameID string) (*MoveResponse, error) {
	game, err := gs.getGame(gameID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch game")
	}

	win, err := game.MakeMove(request.Row, request.Col)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make move")
	}

	if win {
		if err = gs.gameRepository.SaveResult(context.Background(), game.Player1, gameID, "win"); err != nil {
			return nil, errors.Wrap(err, "failed to save game result")
		}
	}

	if !game.Table.IsFull() {
		row, col := game.Table.FindBestMove()
		winRobot, errBot := game.MakeMove(row, col)
		if errBot != nil {
			return nil, errors.Wrap(err, "failed to make move")
		}

		if winRobot {
			err = gs.gameRepository.SaveResult(context.Background(), game.Player1, gameID, "lose")
			if err != nil {
				return nil, errors.Wrap(err, "failed to save game result")
			}
			response := &MoveResponse{
				Board:  game.GetBoard(),
				Result: "robot",
			}
			return response, nil
		}
	} else { //draw
		err = gs.gameRepository.SaveResult(context.Background(), game.Player1, gameID, "draw")
		if err != nil {
			return nil, errors.Wrap(err, "failed to save game result")
		}
		response := &MoveResponse{
			Board:  game.GetBoard(),
			Result: "draw",
		}
		return response, nil
	}

	response := &MoveResponse{
		Board:  game.GetBoard(),
		Result: "",
	}

	return response, nil
}

func (gs *ServiceImpl) getGame(gameID string) (*domain.Game, error) {
	value, ok := gs.games.Load(gameID)
	if !ok {
		return nil, errors.New("game not found")
	}

	return value.(*domain.Game), nil
}

func (gs *ServiceImpl) GetUserInfo(ctx context.Context, userID string) (*UserInfoResponse, error) {
	login, created, err := gs.userRepository.GetUserInfoByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	response := &UserInfoResponse{
		Login:   login,
		Created: created,
	}

	return response, nil
}

func (gs *ServiceImpl) GetCurrentGames() *CurrentGamesResponse {
	var games []string
	gs.games.Range(func(key, _ any) bool {
		if id, ok := key.(string); ok {
			games = append(games, id)
		}
		return true
	})

	response := &CurrentGamesResponse{
		Games: games,
	}

	return response
}

func (gs *ServiceImpl) GetGameInfo(gameID string) (*GameInfoResponse, error) {
	game, err := gs.getGame(gameID)
	if err != nil {
		return nil, errors.New("game not found")
	}

	response := &GameInfoResponse{
		Board: game.GetBoard(),
		Mode:  game.Mode,
	}

	return response, nil
}

func (gs *ServiceImpl) GetFinishedGames(ctx context.Context, userUUID string) (*FinishedGamesResponse, error) {
	games, err := gs.gameRepository.GetFinishedGames(ctx, userUUID)
	if err != nil {
		return nil, errors.New("game not found")
	}

	response := &FinishedGamesResponse{
		FinishedGames: games,
	}

	return response, nil
}

func (gs *ServiceImpl) GetTopPlayers(ctx context.Context, top int) (*TopPlayersResponse, error) {
	topPlayers, err := gs.userRepository.GetTopPlayersWithLogin(ctx, top)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch top players")
	}

	response := &TopPlayersResponse{
		Result: topPlayers,
	}

	return response, nil
}
