package game

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Button struct {
    widget.Button

    LeftClick func()
    RightClick func()
}

func (b *Button) Tapped(_ *fyne.PointEvent) {
    b.LeftClick()
}

func (b *Button) TappedSecondary(_ *fyne.PointEvent) {
    b.RightClick() 
}

func NewButton() *Button {
    b := &Button{}
    b.ExtendBaseWidget(b)

    return b
}
