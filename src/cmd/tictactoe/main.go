package main

import (
	"go.uber.org/fx"
	"net/http"
	"src/internal/delivery/server"
)

func main() {
	app := fx.New(
		fx.Provide(
			server.NewLogger,
			server.NewGameHandler,
			server.NewMux,
			server.NewServer,
		),
		fx.Invoke(func(*http.Server) {}),
	)

	app.Run()
}
