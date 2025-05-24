package main

import (
	"Go_TicTacToe/internal/di"
	"Go_TicTacToe/internal/storage"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	app := fx.New(
		di.DatabaseModule(),
		di.AuthModule(),
		di.GamesModule(),
		di.CoreModule(),

		fx.Invoke(
			storage.RunMigrations,
			func(*http.Server) {},
		),
	)

	app.Run()
}
