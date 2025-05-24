package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type GameRepository interface {
	SaveNewGame(ctx context.Context, gameID, mode string) error
	SaveResult(ctx context.Context, userID, gameID, result string) error
	GetFinishedGames(ctx context.Context, userID string) ([]string, error)
}

type GameRepositoryImpl struct {
	db *pgxpool.Pool
}

var _ GameRepository = (*GameRepositoryImpl)(nil)

func NewGameRepository(db *pgxpool.Pool) GameRepository {
	return &GameRepositoryImpl{db: db}
}

const (
	saveNewGameQuery = `
	INSERT INTO games (id, mode) VALUES ($1, $2)
	`

	saveResultQuery = `
	INSERT INTO game_results(user_id, game_id, result) VALUES ($1, $2, $3)
	`

	getFinishedGameQuery = `
    SELECT DISTINCT game_id FROM game_results
    WHERE user_id = $1 AND (result = 'win' OR result = 'draw')
    `

	getWinGameQuery = `
    SELECT  user_id,
        	COUNT(game_id) AS win_cnt
    FROM game_results
    WHERE result = 'win'
    GROUP BY user_id`

	getTotalGamesQuery = `
	SELECT  user_id,
	    	COUNT(game_id) AS game_cnt
	FROM game_results
	WHERE result = 'lose' OR result = 'draw'
	GROUP BY user_id`
)

func (gr *GameRepositoryImpl) SaveNewGame(ctx context.Context, gameID, mode string) error {

	_, err := gr.db.Exec(ctx, saveNewGameQuery, gameID, mode)
	if err != nil {
		return errors.Wrap(err, "failed to save game")
	}

	return nil
}

func (gr *GameRepositoryImpl) SaveResult(ctx context.Context, userID, gameID, result string) error {
	_, err := gr.db.Exec(ctx, saveResultQuery, userID, gameID, result)
	if err != nil {
		return errors.Wrap(err, "failed to save game result")
	}

	return nil
}

func (gr *GameRepositoryImpl) GetFinishedGames(ctx context.Context, userID string) ([]string, error) {
	var games []string
	rows, err := gr.db.Query(ctx, getFinishedGameQuery, userID)
	if err != nil {
		return games, errors.Wrap(err, "failed to query finished games")
	}

	for rows.Next() {
		var gameID string
		if err = rows.Scan(&gameID); err != nil {
			return games, errors.Wrap(err, "failed to query game ID")
		}

		games = append(games, gameID)
	}

	if err = rows.Err(); err != nil {
		return games, errors.Wrap(err, "failed to scan finished games")
	}

	return games, nil
}
