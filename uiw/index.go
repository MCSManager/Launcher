package uiw

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type MyLabel struct {
	Canvas *canvas.Text
}

func NewMyLabel(text string) *MyLabel {
	obj := canvas.NewText(text, color.Black)
	obj.Resize(fyne.NewSize(100, 200))
	obj.Refresh()
	// obj.Refresh()
	return &MyLabel{Canvas: obj}
}

func (p *MyLabel) SetColor(color color.Color) {
	p.Canvas.Color = color
	p.Canvas.Refresh()
}

func (p *MyLabel) SetText(text string) {
	p.Canvas.Text = text
	p.Canvas.Refresh()
}
func (p *MyLabel) SetFontSize(size float32) {
	p.Canvas.TextSize = size
	p.Canvas.Refresh()
}
