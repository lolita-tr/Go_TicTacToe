package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, userID, login, password string) error
	GetUserByLogin(ctx context.Context, login string) (string, string, error)
	GetUserUUIDbyLogin(ctx context.Context, login string) (string, error)
	GetUserInfoByID(ctx context.Context, userID string) (string, time.Time, error)
	GetTopPlayers(ctx context.Context, top int) (map[string]float64, error)
	GetTopPlayersWithLogin(ctx context.Context, top int) ([]PlayerWinRate, error)
	GetFinishedGamesByToken(ctx context.Context, userID string) ([]string, error)
}

type UserRepositoryImpl struct {
	db *pgxpool.Pool
}

var _ UserRepository = (*UserRepositoryImpl)(nil)

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &UserRepositoryImpl{db: db}
}

type PlayerWinRate struct {
	UserUUID string
	Login    string
	WinRate  float64
}

const (
	signInQuery = `
	INSERT INTO users (id, login, password) VALUES ($1, $2, $3)
`
	getUserQuery = `
	SELECT id, password FROM users WHERE login = $1
`
	getUserIdQuery = `
	SELECT id FROM users WHERE login = $1
`
	getUserInfoQuery = `
	SELECT login, created_at FROM users WHERE id=$1
`
	getFinishedGameBYTokenQuery = `
    SELECT DISTINCT game_id FROM game_results
    WHERE user_id = $1 AND (result = 'win' OR result = 'draw')
    `

	getWinRateQuery = `
	WITH win_games AS (
	SELECT  user_id,
       	COUNT(game_id) AS win_cnt
    FROM game_results
    WHERE result = 'win'
    GROUP BY user_id),
	
	total_games AS (
	SELECT  user_id,
	   	COUNT(game_id) AS game_cnt
	FROM game_results
	WHERE result = 'lose' OR result = 'draw'
	GROUP BY user_id
	)
	
	SELECT total_games.user_id,
			ROUND(win_games.win_cnt::decimal/total_games.game_cnt, 2) AS win_rate
    FROM total_games
	INNER JOIN win_games ON total_games.user_id = win_games.user_id
	ORDER BY win_rate DESC
	LIMIT $1 
	`

	getWinRateLoginQuery = `
	WITH win_games AS (
	SELECT  user_id,
	COUNT(game_id) AS win_cnt
	FROM game_results
	WHERE result = 'win'
	GROUP BY user_id),
	
	total_games AS (
	SELECT  user_id,
	COUNT(game_id) AS game_cnt
	FROM game_results
	WHERE result = 'lose' OR result = 'draw'
	GROUP BY user_id
	)
	
	SELECT total_games.user_id,
	users.login,
	ROUND(win_games.win_cnt::decimal/total_games.game_cnt, 2) AS win_rate
	FROM total_games
	INNER JOIN win_games ON total_games.user_id = win_games.user_id
	INNER JOIN users ON total_games.user_id = users.id
	ORDER BY win_rate DESC
	LIMIT $1
`
)

func (ur *UserRepositoryImpl) CreateUser(ctx context.Context, userID, login, password string) error {
	_, err := ur.db.Exec(ctx, signInQuery, userID, login, password)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	return nil
}

func (ur *UserRepositoryImpl) GetUserByLogin(ctx context.Context, login string) (string, string, error) {
	var userID, hashedPassword string
	err := ur.db.QueryRow(ctx, getUserQuery, login).Scan(&userID, &hashedPassword)
	if err != nil {
		return "", "", errors.Wrap(err, "user not found")
	}

	return userID, hashedPassword, nil
}

func (ur *UserRepositoryImpl) GetUserUUIDbyLogin(ctx context.Context, login string) (string, error) {
	var userID string
	err := ur.db.QueryRow(ctx, getUserIdQuery, login).Scan(&userID)
	if err != nil {
		return "", errors.Wrap(err, "user not found")
	}

	return userID, nil
}

func (ur *UserRepositoryImpl) GetUserInfoByID(ctx context.Context, userID string) (string, time.Time, error) {
	var login string
	var created time.Time

	err := ur.db.QueryRow(ctx, getUserInfoQuery, userID).Scan(&login, &created)
	if err != nil {
		return "", time.Time{}, errors.Wrap(err, "user not found")
	}

	return login, created, nil
}

func (gr *UserRepositoryImpl) GetTopPlayers(ctx context.Context, top int) (map[string]float64, error) {
	topPlayers := make(map[string]float64)

	rows, err := gr.db.Query(ctx, getWinRateQuery, top)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query top players")
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		var winRate float64
		if err = rows.Scan(&userID, &winRate); err != nil {
			return nil, errors.Wrap(err, "failed to scan top players")
		}

		topPlayers[userID] = winRate
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to scan top players")
	}

	return topPlayers, nil
}

func (gr *UserRepositoryImpl) GetTopPlayersWithLogin(ctx context.Context, top int) ([]PlayerWinRate, error) {
	rows, err := gr.db.Query(ctx, getWinRateLoginQuery, top)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get top players")
	}

	defer rows.Close()

	var result []PlayerWinRate
	for rows.Next() {
		var p PlayerWinRate
		if err = rows.Scan(&p.UserUUID, &p.Login, &p.WinRate); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		result = append(result, p)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows error")
	}

	return result, nil
}

func (gr *UserRepositoryImpl) GetFinishedGamesByToken(ctx context.Context, userID string) ([]string, error) {
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
