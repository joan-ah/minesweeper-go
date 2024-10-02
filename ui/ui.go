package ui

import (
   "fyne.io/fyne/v2"
   "fyne.io/fyne/v2/container"

   "minesweeper/game"
)

var c *fyne.Container

func Render(window fyne.Window, game *game.Game) {
   c = container.NewGridWithColumns(game.Cols)
   
   for _, cell := range game.Cells {
      c.Add(cell.Container)
   }

   w := float32(game.Cols * 50)
   h := float32(game.Rows * 50)

   window.Resize(fyne.NewSize(w, h))

   window.SetContent(container.NewBorder(
      nil,
      nil,
      nil,
      nil,
      c,
   ))
}
