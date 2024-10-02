package main

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"

    "minesweeper/game"
)

func main() {
    app := app.New()
    window := app.NewWindow("Minesweeper")
    if icon, err := fyne.LoadResourceFromPath("./bomb.png"); err == nil {
	window.SetIcon(icon)
    }

    game := game.New(10, 10, "easy")

    game.Render()
    window.ShowAndRun()
}
