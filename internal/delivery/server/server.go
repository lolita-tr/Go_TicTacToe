package server

import (
	"Go_TicTacToe/internal/auth"
	"Go_TicTacToe/internal/game"
	"Go_TicTacToe/internal/middleware"
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

func NewMux(gameHandler *game.Handler, authHandler *auth.Handler, auth *middleware.UserAuthenticator) *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/", fileServer)

	//публичные поинты
	mux.HandleFunc("/signup", authHandler.SignUp)
	mux.HandleFunc("/login", authHandler.Authorization)
	mux.HandleFunc("/refresh/access", authHandler.RefreshAccessToken)
	mux.HandleFunc("/refresh/refresh", authHandler.RefreshToken)

	//с авторизацией
	mux.HandleFunc("/new/bot", auth.AuthMiddleWare(gameHandler.CreateBotGame))
	mux.HandleFunc("/new/user", auth.AuthMiddleWare(gameHandler.CreateUserGame))
	mux.HandleFunc("/move/bot", auth.AuthMiddleWare(gameHandler.MakeMoveBot))
	mux.HandleFunc("/move/users", auth.AuthMiddleWare(gameHandler.MakeMovePlayers))
	mux.HandleFunc("/board", auth.AuthMiddleWare(gameHandler.GetBoard))
	mux.HandleFunc("/reset", auth.AuthMiddleWare(gameHandler.ResetGame))

	mux.HandleFunc("/userinfo/id", auth.AuthMiddleWare(gameHandler.GetUserInfo))
	mux.HandleFunc("/game", auth.AuthMiddleWare(gameHandler.GetGame))
	mux.HandleFunc("/games/current", auth.AuthMiddleWare(gameHandler.GetCurrentGames))
	mux.HandleFunc("/games/finished/id", auth.AuthMiddleWare(gameHandler.GetFinishedGames))
	mux.HandleFunc("/top/players", auth.AuthMiddleWare(gameHandler.GetTopPlayers))

	return mux
}

func NewServer(lc fx.Lifecycle, mux *http.ServeMux, logger *zap.Logger) *http.Server {
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting Tic Tac Toe server on :8080")
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down TIc Tac Toe server")
			return server.Shutdown(ctx)
		},
	})
	return server
}

func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}
