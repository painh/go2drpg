package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type Cursor struct {
	x                float64
	y                float64
	showCursor       bool
	prompt           string
	lastCursorStatus bool
}

func (c *Cursor) Init() {
	c.x = float64(ConfigInstance.CursorX)
	c.y = float64(ConfigInstance.CursorY)

}

func (c *Cursor) MakeTimePrompt(a ...interface{}) string {
	text := fmt.Sprint(a...)

	hour := GameInstance.player.curTimeMin / 60
	min := GameInstance.player.curTimeMin % 60
	return fmt.Sprintf("%d/%02d:%02d> %s", GameInstance.player.day, int(hour), int(min), text)
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
