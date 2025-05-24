package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(p DBParams) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), p.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	if err1 := dbpool.Ping(context.Background()); err1 != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err1)
	}

	return dbpool, nil
}
