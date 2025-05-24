package domain

import (
	"math"
)

const (
	TABLE_SIZE = 3
	WIN_POINTS = 3

	PLAYER_X = 1
	PLAYER_O = 2
)

type Table struct {
	Table [][]int
}

func createTable() *Table {
	table := &Table{
		Table: make([][]int, TABLE_SIZE),
	}
	for i := range table.Table {
		table.Table[i] = make([]int, TABLE_SIZE)
	}
	return table
}

func (t *Table) FindBestMove() (int, int) {
	bestScore := math.Inf(-1)
	bestMove := [2]int{-1, -1}

	for i := 0; i < TABLE_SIZE; i++ {
		for j := 0; j < TABLE_SIZE; j++ {
			if t.Table[i][j] == 0 {
				t.Table[i][j] = PLAYER_X
				score := t.minimax(0, false)
				t.Table[i][j] = 0

				if score > bestScore {
					bestScore = score
					bestMove = [2]int{i, j}
				}
			}
		}
	}
	return bestMove[0], bestMove[1]
}

func (t *Table) minimax(depth int, isMax bool) float64 {
	if t.checkWinAI(PLAYER_X) {
		return 10 - float64(depth)
	}

	if t.checkWinAI(PLAYER_O) {
		return -10 + float64(depth)
	}

	if t.IsFull() {
		return 0
	}

	if isMax {
		bestScore := math.Inf(-1)
		for i := 0; i < TABLE_SIZE; i++ {
			for j := 0; j < TABLE_SIZE; j++ {
				if t.Table[i][j] == 0 {
					t.Table[i][j] = PLAYER_X
					score := t.minimax(depth+1, false)
					t.Table[i][j] = 0
					bestScore = math.Max(bestScore, score)
				}
			}
		}
		return bestScore
	} else {
		bestScore := math.Inf(1)
		for i := 0; i < TABLE_SIZE; i++ {
			for j := 0; j < TABLE_SIZE; j++ {
				if t.Table[i][j] == 0 {
					t.Table[i][j] = PLAYER_O
					score := t.minimax(depth+1, true)
					t.Table[i][j] = 0
					bestScore = math.Min(bestScore, score)
				}
			}
		}
		return bestScore
	}
}

func (t *Table) IsFull() bool {
	for i := 0; i < TABLE_SIZE; i++ {
		for j := 0; j < TABLE_SIZE; j++ {
			if t.Table[i][j] == 0 {
				return false
			}
		}
	}
	return true
}

func (t *Table) checkWin(x, y int) bool {
	return t.checkRow(x, y) || t.checkColumn(x, y) || t.checkDiagonal1(x, y) || t.checkDiagonal2(x, y)
}

func (t *Table) checkWinAI(player int) bool {
	for i := 0; i < TABLE_SIZE; i++ {
		for j := 0; j < TABLE_SIZE; j++ {
			if t.Table[i][j] == player && (t.checkRow(i, j) || t.checkColumn(i, j) || t.checkDiagonal1(i, j) || t.checkDiagonal2(i, j)) {
				return true
			}
		}
	}
	return false
}

func (t *Table) checkRow(i, j int) bool {
	d1 := 0
	for d1 <= j && t.Table[i][j-d1] == t.Table[i][j] {
		d1++
	}
	d2 := 0
	for j+d2 < TABLE_SIZE && t.Table[i][j+d2] == t.Table[i][j] {
		d2++
	}
	return d1+d2 > WIN_POINTS
}

func (t *Table) checkColumn(i, j int) bool {
	d1 := 0
	for d1 <= i && t.Table[i-d1][j] == t.Table[i][j] {
		d1++
	}
	d2 := 0
	for i+d2 < TABLE_SIZE && t.Table[i+d2][j] == t.Table[i][j] {
		d2++
	}
	return d1+d2 > WIN_POINTS
}

func (t *Table) checkDiagonal1(i, j int) bool {
	d1 := 0
	for d1 <= j && d1 <= i && t.Table[i-d1][j-d1] == t.Table[i][j] {
		d1++
	}
	d2 := 0
	for i+d2 < TABLE_SIZE && j+d2 < TABLE_SIZE && t.Table[i+d2][j+d2] == t.Table[i][j] {
		d2++
	}
	return d1+d2 > WIN_POINTS
}

func (t *Table) checkDiagonal2(i, j int) bool {
	d1 := 0
	for i+d1 < TABLE_SIZE && d1 <= j && t.Table[i+d1][j-d1] == t.Table[i][j] {
		d1++
	}
	d2 := 0
	for d2 <= i && j+d2 < TABLE_SIZE && t.Table[i-d2][j+d2] == t.Table[i][j] {
		d2++
	}
	return d1+d2 > WIN_POINTS
}
