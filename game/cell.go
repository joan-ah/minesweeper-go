package game

import (
    "fmt"
    "image/color"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"
)

type Cell struct {
    Container *fyne.Container

    game *Game

    X, Y int
    surroundingMines *int
    Mine, Flagged, Revealed bool
}

func (cell *Cell) CreateContainer() {
    if cell.Container != nil {
	cell.Container.Objects[0].(*canvas.Rectangle).FillColor = color.NRGBA{ 220, 220, 220, 255 }
	b := cell.Container.Objects[1].(*Button)
	b.SetText("")
	b.Enable()
	return
    }

    r := canvas.NewRectangle(color.NRGBA{ 220, 220, 220, 255 })

    b := NewButton()
    b.LeftClick = func () { cell.Reveal(true) }
    b.RightClick = cell.Flag
    b.Importance = widget.LowImportance

    c := container.NewStack(r, b)
    
    cell.Container = c
}

func (cell *Cell) Reveal(direct bool) {
    if cell.Flagged || cell.game.GameOver {
	return
    }

    if cell.Revealed {
	if direct {
	    cell.RevealSurrounding()
	}
	return
    }

    cell.Revealed = true

    if cell.Mine {
	cell.Container.Objects[0].(*canvas.Rectangle).FillColor = color.NRGBA{ 255, 0, 0, 200 }
	cell.game.GameOver = true

	dialog.ShowConfirm(
	    "Game Over!",
	    "You exploded :(",
	    func (_ bool) { cell.game.Restart() },
	    fyne.CurrentApp().Driver().AllWindows()[0],
	)

	return
    }

    if cell.game.Victory() {
	cell.game.GameOver = true

	dialog.ShowConfirm(
	    "Game Over!",
	    "You win :D",
	    func (_ bool) { cell.game.Restart() },
	    fyne.CurrentApp().Driver().AllWindows()[0],
	)

	return

    }
    
    mines := cell.SurroundingMines()

    cell.Container.Objects[0].(*canvas.Rectangle).FillColor = color.Black

    if mines > 0 {
	cell.Container.Objects[1].(*Button).SetText(fmt.Sprintf("%d", mines))

	return
    }

    cell.Container.Objects[1].(*Button).Disable()

    cell.RevealSurrounding()
}

func (cell *Cell) Flag() {
    if cell.Revealed || cell.game.GameOver {
	return
    }

    var c color.Color

    if cell.Flagged {
	c = color.NRGBA{ 220, 220, 220, 255 }
    } else {
	c = color.NRGBA{ 135, 206, 235, 255 }
    }

    cell.Flagged = !cell.Flagged
    cell.Container.Objects[0].(*canvas.Rectangle).FillColor = c
}

func (cell *Cell) SurroundingMines() int {
    if cell.surroundingMines != nil {
	return *cell.surroundingMines
    }

    mines := 0

    for y := -1; y <= 1; y++ {
	for x := -1; x <= 1; x++ {
	    check_x := cell.X + x
	    check_y := cell.Y + y

	    if (x == 0 && y == 0) || !cell.game.WithinBounds(check_x, check_y) {
		continue
	    }

	    if cell.game.CellAt(check_x, check_y).Mine {
		mines++
	    }
	}
    }

    cell.surroundingMines = &mines 

    return mines
}

func (cell *Cell) RevealSurrounding() {
    for y := -1; y <= 1; y++ {
	for x := -1; x <= 1; x++ {
	    check_x := cell.X + x
	    check_y := cell.Y + y

	    if (x == 0 && y == 0) || !cell.game.WithinBounds(check_x, check_y) {
		continue
	    }

	    cell.game.CellAt(check_x, check_y).Reveal(false)
	}
    }
}

