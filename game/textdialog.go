package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"strings"
)

type TextDialog struct {
	x     float64
	y     float64
	lines []string
}

func (t *TextDialog) SetText(text string) {
	t.lines = strings.Split(text, "\n")
}

func (t *TextDialog) Update() {
}

func (t *TextDialog) Draw(screen *ebiten.Image) {
	y := t.y
	for _, v := range t.lines {
		ret := defaultFontInstance.DrawTextInBox(screen, v, t.x, y)
		y += ret
	}
}
