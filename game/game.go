package game

import (
    "fmt"
    "math/rand"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/data/binding"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"
)

type Game struct {
    c *fyne.Container

    Rows, Cols, Mines int
    Difficulty string
    GameOver bool
    Flags binding.Int

    Cells []Cell
}

func New(rows, cols int, difficulty string) Game {
    mines := calculateNumMinesForDifficulty(rows, cols, difficulty)

    game := Game{
	nil,
	rows,
	cols,
	mines,
	difficulty,
	false,
	binding.NewInt(),
	make([]Cell, rows * cols),
    }

    if mines > game.MaxMines() {
	game.Mines = game.MaxMines()
    }

    game.Restart()

    return game
}

func (game *Game) Restart() {
    game.GameOver = false
    game.Flags.Set(0)

    for y := 0; y < game.Rows; y++ {
	for x := 0; x < game.Cols; x++ {
	    cell := game.CellAt(x, y)
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

    return revealed == g.Rows * g.Cols - g.Mines
}

func (g *Game) CellAt(x, y int) *Cell {
    return &g.Cells[y * g.Cols + x]
}

func (g *Game) WithinBounds(x, y int) bool {
    return x >= 0 && x < g.Cols && y >= 0 && y < g.Rows
}

func (g *Game) MaxMines() int {
    return g.Rows * g.Cols - 1
}

func (game *Game) Render() {
    window := defaultWindow()
    game.c = container.NewGridWithColumns(game.Cols)
   
    for y := 0; y < game.Rows; y++ {
	for x := 0; x < game.Cols; x++ {
	    game.c.Add(game.CellAt(x, y).Container)
	}
    }

    w := float32(game.Cols * 50)
    h := float32(game.Rows * 50 + 50)

    window.Resize(fyne.NewSize(w, h))

    minesLeft := binding.NewString()

    game.Flags.AddListener(binding.NewDataListener(func() {
	flags, _ := game.Flags.Get()
	minesLeft.Set(fmt.Sprintf("Mines left: %d", game.Mines - flags)) 
    }))

    top := container.NewHBox(
	widget.NewLabel(fmt.Sprintf("Size: %d x %d", game.Cols, game.Rows)),
	widget.NewButtonWithIcon("", theme.SettingsIcon(), game.ShowSettings), 
	layout.NewSpacer(),
	widget.NewLabelWithData(minesLeft),
	widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), game.Restart),
    )

    window.SetContent(container.NewBorder(
	top,
	nil,
	nil,
	nil,
	game.c,
    ))
}

func (game *Game) ShowSettings() {
    var difficulty string
    diff := widget.NewRadioGroup([]string{
	"easy",
	"normal",
	"hard",
	"expert",
    },
    func(s string) {
	difficulty = s
    })
    diff.SetSelected(game.Difficulty)

    rowsL := widget.NewLabel(fmt.Sprintf("%d", game.Rows))
    rows := widget.NewSlider(10, 50)
    rows.Step = 1
    rows.Value = float64(game.Rows)
    rows.OnChanged = func(f float64) {
	rowsL.SetText(fmt.Sprintf("%d", int64(rows.Value)))	
    }

    colsL := widget.NewLabel(fmt.Sprintf("%d", game.Cols))
    cols := widget.NewSlider(10, 50)
    cols.Step = 1
    cols.Value = float64(game.Cols)
    cols.OnChanged = func(f float64) {
	colsL.SetText(fmt.Sprintf("%d", int64(cols.Value)))	
    }

    dialog.ShowCustomConfirm(
	"Game Settings",
	"Update",
	"Cancel",
	container.NewStack(container.NewVBox(
	    container.NewBorder(nil, nil, widget.NewLabel("Rows"), rowsL, rows),
	    container.NewBorder(nil, nil, widget.NewLabel("Cols"), colsL, cols),
	    container.NewBorder(widget.NewLabel("Mines"), nil, nil, nil, diff),
	)),
	func(b bool) {
	    if !b {
		return
	    }
	    
	    *game = New(int(rows.Value), int(cols.Value), difficulty)
	    game.Render()
	},
	defaultWindow(),
    )
}

func calculateNumMinesForDifficulty(rows, cols int, difficulty string) int {
    var ratio float64

    switch difficulty {
    case "easy":
	ratio = 0.15
	break
    case "normal":
	ratio = 0.18
	break
    case "hard":
	ratio = 0.20
	break
    case "expert":
	ratio = 0.25
	break
    default:
	ratio = 0.99
    }

    return int(float64(rows * cols) * ratio)
}

func defaultWindow() fyne.Window {
    return fyne.CurrentApp().Driver().AllWindows()[0]
}
