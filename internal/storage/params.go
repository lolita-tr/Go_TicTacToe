package storage

import (
	"fmt"
	"os"
)

type DBParams struct {
	URL string
}

func NewPostgresParams() DBParams {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	database := os.Getenv("DB_DATABASE")

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, database)
	return DBParams{URL: url}
}
