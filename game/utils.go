package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"time"
)

func DrawRect(screen *ebiten.Image, x, y, width, height float64, color color.Color) {
	width--
	height--
	ebitenutil.DrawLine(screen, x, y, x+width, y, color)
	ebitenutil.DrawLine(screen, x+width, y, x+width, y+height, color)
	ebitenutil.DrawLine(screen, x, y+height, x+width, y+height, color)
	ebitenutil.DrawLine(screen, x, y, x, y+height, color)
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
