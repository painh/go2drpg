package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
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
	c.x = float64(SettingConfigInstance.CursorX)
	c.y = float64(SettingConfigInstance.CursorY)

}

func (c *Cursor) MakeTimePrompt(a ...interface{}) string {
	text := fmt.Sprint(a...)

	time := GameInstance.player.GetTimeString()
	return fmt.Sprintf("%s> %s", time, text)
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
	defaultFontInstance.Draw(dst, c.prompt, SettingConfigInstance.CursorX, SettingConfigInstance.CursorY, color.RGBA{255, 255, 255, 255})
}
