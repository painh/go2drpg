package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type Cursor struct {
	x                float64
	y                float64
	day              int
	curTimeMin       int
	showCursor       bool
	prompt           string
	lastCursorStatus bool
}

func (c *Cursor) Init() {
	c.x = float64(ConfigInstance.CursorX)
	c.y = float64(ConfigInstance.CursorY)
	c.curTimeMin = ConfigInstance.StartTimeMin
}

func (c *Cursor) MakeTimePrompt(a ...interface{}) string {
	text := fmt.Sprint(a...)

	hour := c.curTimeMin / 60
	min := c.curTimeMin % 60
	return fmt.Sprintf("%d/%02d:%02d> %s", c.day, int(hour), int(min), text)
}

func (c *Cursor) Update() {
	c.showCursor = math.Sin(float64(makeTimestamp()/100)) > 0
	if c.showCursor == c.lastCursorStatus {
		return
	}

	c.lastCursorStatus = c.showCursor

	cursor := ""
	if c.showCursor {
		cursor = "_"
	}

	c.prompt = c.MakeTimePrompt(cursor)
}

func (c *Cursor) Draw(dst *ebiten.Image) {
	defaultFontInstance.Draw(dst, c.prompt, ConfigInstance.CursorX, ConfigInstance.CursorY)
}
