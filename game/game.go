package game

import "math/rand"

type Game struct {
    Rows, Cols int
    Mines uint
    GameOver bool

    Cells []Cell
}

func New(rows, cols int, mines uint) Game {
    game := Game{ rows, cols, mines, false, make([]Cell, rows * cols) }
    game.Restart()

    return game
}

func (game *Game) Restart() {
    game.GameOver = false

    for y := 0; y < game.Cols; y++ {
	for x := 0; x < game.Rows; x++ {
	    cell := &game.Cells[y * game.Rows + x]
	    cell.game = game
	    cell.Mine = false
	    cell.Flagged = false
	    cell.Revealed = false
	    cell.surroundingMines = nil
	    cell.X = x
	    cell.Y = y
	    cell.CreateContainer()
	}
    }

    mines := game.Mines

    for mines > 0 {
	x := rand.Int() % game.Cols 
	y := rand.Int() % game.Rows

	cell := game.CellAt(x, y)

	if cell.Mine {
	    continue
	}

	cell.Mine = true
	mines--
    }

}

func (g *Game) Victory() bool {
    revealed := 0

    for y := 0; y < g.Rows; y++ {
	for x := 0; x < g.Cols; x++ {
	    if g.CellAt(x, y).Revealed {
		revealed++
	    }
	}
    }

    return revealed == g.Rows * g.Cols - int(g.Mines)
}

func (g *Game) CellAt(x, y int) *Cell {
    return &g.Cells[y * g.Rows + x]
}

func (g *Game) WithinBounds(x, y int) bool {
    return x >= 0 && x < g.Cols && y >= 0 && y < g.Rows
}
