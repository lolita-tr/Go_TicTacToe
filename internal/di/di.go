package di

import (
	"Go_TicTacToe/internal/auth"
	"Go_TicTacToe/internal/delivery/server"
	"Go_TicTacToe/internal/game"
	"Go_TicTacToe/internal/middleware"
	"Go_TicTacToe/internal/storage"
	"go.uber.org/fx"
)

func CoreModule() fx.Option {
	return fx.Provide(
		server.NewLogger,
		server.NewMux,
		server.NewServer)

}

func GamesModule() fx.Option {
	return fx.Provide(
		game.NewGameService,
		game.NewGameHandler,
	)
}

func AuthModule() fx.Option {
	return fx.Provide(
		auth.NewJwtProvider,
		auth.NewAuthService,
		auth.NewAuthHandler,
		middleware.NewUserAuthenticator)
}

func DatabaseModule() fx.Option {
	return fx.Provide(
		storage.NewPostgresParams,
		storage.NewPostgresDB,
		storage.NewUserRepository,
		storage.NewGameRepository,
	)
}
