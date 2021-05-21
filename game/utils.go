package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/k0kubun/pp"
	"image/color"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

func DrawRect(screen *ebiten.Image, x, y, width, height float64, color color.Color) {
	width--
	height--
	ebitenutil.DrawLine(screen, x, y, x+width, y, color)                 // top
	ebitenutil.DrawLine(screen, x+width, y, x+width, y+height, color)    // right
	ebitenutil.DrawLine(screen, x, y, x, y+height, color)                // left
	ebitenutil.DrawLine(screen, x-1, y+height, x+width, y+height, color) // bottom
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func ManhattanDistance(x1, y1, x2, y2 float64) float64 {
	return math.Abs(x1-x2) + math.Abs(y1-y2)
}

func atoi(s string) int {
	i, err := strconv.Atoi(strings.TrimSpace(s))

	if err != nil {
		log.Fatalln(err)
	}

	return i
}

func dump(data interface{}) {
	pp.Print(data)
}
