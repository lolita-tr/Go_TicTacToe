package domain

import (
	"errors"
	"sync"
)

type Game struct {
	Table  *Table
	Player int
	mu     sync.Mutex
}

var (
	errCellTaken = errors.New("cell is already taken")
)

func CreateGame() *Game {
	return &Game{
		Table:  createTable(),
		Player: PLAYER_O,
	}
}

func (g *Game) GetBoard() [][]int {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.Table.Table
}

func (g *Game) CurrentPlayer() int {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.Player
}

func (g *Game) MakeMove(row, col int) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Table.Table[row][col] != 0 {
		return false, errCellTaken
	}

	g.Table.Table[row][col] = g.Player

	if g.Table.checkWin(row, col) {
		return true, nil
	}

	if g.Player == PLAYER_O {
		g.Player = PLAYER_X
	} else {
		g.Player = PLAYER_O
	}

	return false, nil
}

func (g *Game) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.Table = createTable()
	g.Player = PLAYER_O
}
