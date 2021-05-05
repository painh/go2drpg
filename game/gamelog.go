package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"strings"
)

type GameLogElement struct {
	text       string
	createdTs  int64
	imageBuf   *ebiten.Image
	lineHeight int
}

func (g *GameLogElement) Set(text string) {
	g.text = text

	newText := ""
	resultText := ""
	linesCnt := 1
	y := 0
	for _, v := range text {
		prevText := newText
		newText += string(v)
		rect := defaultFontInstance.BoundString(newText)
		width := rect.Dx()
		if width >= ConfigInstance.LogWidth {
			rect := defaultFontInstance.BoundString(prevText)
			y += rect.Dy() + ConfigInstance.LineSpacing

			resultText = prevText + "\n"
			prevText = newText
			newText = string(v)
			linesCnt++
		}
	}
	resultText += newText
	rect := defaultFontInstance.BoundString(newText)
	y += rect.Dy() + ConfigInstance.LineSpacing
	y += ConfigInstance.LineSpacing

	strs := strings.Split(resultText, "\n")

	g.lineHeight = y
	g.imageBuf = ebiten.NewImage(ConfigInstance.LogWidth, g.lineHeight)

	y = 0
	for _, v := range strs {
		rect := defaultFontInstance.BoundString(v)
		y += rect.Dy() + ConfigInstance.LineSpacing
		defaultFontInstance.Draw(g.imageBuf, v, 0, y)
	}

}

func (g *GameLogElement) Update() {

}

func (g *GameLogElement) Draw(screen *ebiten.Image, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	screen.DrawImage(g.imageBuf, op)
}

type GameLog struct {
	lines []GameLogElement
}

var GameLogInstance = GameLog{lines: []GameLogElement{}}

func (g *GameLog) Update() {
	for _, e := range g.lines {
		e.Update()
	}
}

func (g *GameLog) Draw(screen *ebiten.Image) {
	y := float64(GameInstance.screenHeight)

	for i, e := range g.lines {
		if i >= ConfigInstance.LogLines {
			break
		}

		y -= float64(e.lineHeight)
		e.Draw(screen, 0, y)
	}
}

func (g *GameLog) Add(a ...interface{}) {

	text := fmt.Sprint(a...)

	l := GameLogElement{text: text, createdTs: makeTimestamp()}
	l.Set(text)

	g.lines = append([]GameLogElement{l}, g.lines...)

	if len(g.lines) > ConfigInstance.LogLines {
		g.lines = g.lines[:len(g.lines)-1]
	}

	fmt.Println(text)
}
