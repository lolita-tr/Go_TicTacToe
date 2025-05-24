package game

import (
	"Go_TicTacToe/internal/domain"
	"context"
)

type Service interface {
	CreateBotGame(userID string) (*NewGameResponse, error)
	CreateUserGame(userID string, oppLogin string) (*NewGameResponse, error)
	MakeMovePlayers(request MoveRequest, gameID string) (*MoveResponse, error)
	MakeMoveBot(request MoveRequest, gameID string) (*MoveResponse, error)
	getGame(gameID string) (*domain.Game, error)
	GetUserInfo(ctx context.Context, userID string) (*UserInfoResponse, error)
	GetCurrentGames() *CurrentGamesResponse
	GetGameInfo(gameID string) (*GameInfoResponse, error)
	GetFinishedGames(ctx context.Context, userUUID string) (*FinishedGamesResponse, error)
	GetTopPlayers(ctx context.Context, top int) (*TopPlayersResponse, error)
}
