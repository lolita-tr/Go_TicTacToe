package game

import (
	"Go_TicTacToe/internal/storage"
	"time"
)

type NewGameUserRequest struct {
	OpponentLogin string `json:"opponent_login"`
}

type NewGameResponse struct {
	GameUUID string `json:"game_id"`
}

type UserInfoRequest struct {
	UserUUID string `json:"user_uuid"`
}

type UserInfoResponse struct {
	Login   string    `json:"login"`
	Created time.Time `json:"created"`
}

type GameInfoRequest struct {
	GameUUID string `json:"game_uuid"`
}
type CurrentGamesResponse struct {
	Games []string `json:"games"`
}

type GameInfoResponse struct {
	Board [][]int `json:"board"`
	Mode  string  `json:"mode"`
}

type MoveRequest struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type MoveResponse struct {
	Board  [][]int `json:"board"`
	Result string  `json:"winner"`
}

type FinishedGamesRequest struct {
	UserUUID string `json:"user_uuid"`
}

type FinishedGamesResponse struct {
	FinishedGames []string `json:"games"`
}

type TopPlayersRequest struct {
	Top int `json:"top"`
}

type TopPlayersResponse struct {
	Result []storage.PlayerWinRate `json:"result"`
}
