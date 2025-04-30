package server

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

func NewMux(logger *zap.Logger, handler *GameHandler) *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/", fileServer)

	mux.HandleFunc("/new", handler.HandleNewGame)
	mux.HandleFunc("/move", handler.HandleMove)
	mux.HandleFunc("/board", handler.HandleBoard)
	mux.HandleFunc("/reset", handler.HandleReset)
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
