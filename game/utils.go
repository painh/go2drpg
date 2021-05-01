package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math"
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

func ManhattanDistance(x1, y1, x2, y2 float64) float64 {
	return math.Abs(x1-x2) + math.Abs(y1-y2)
}
