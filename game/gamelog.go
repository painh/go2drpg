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
	for _, v := range text {
		prevText := newText
		newText += string(v)
		rect := defaultFontInstance.BoundString(newText)
		width := rect.Dx()
		if width >= ConfigInstance.LogWidth {
			resultText = prevText + "\n"
			prevText = newText
			newText = string(v)
			linesCnt++
		}
	}
	resultText += newText

	strs := strings.Split(resultText, "\n")
	y := 0
	const margin = 5
	for _, v := range strs {
		rect := defaultFontInstance.BoundString(v)
		y += rect.Dy() + margin
	}
	y += margin

	//rect := defaultFontInstance.BoundString(resultText)

	//g.lineHeight = int(rect.Dy())
	//g.imageBuf = ebiten.NewImage(ConfigInstance.LogWidth, int(g.lineHeight)+1)
	g.lineHeight = y
	g.imageBuf = ebiten.NewImage(ConfigInstance.LogWidth, g.lineHeight)
	//DrawRect(g.imageBuf, 0, 0, float64(ConfigInstance.LogWidth), float64(g.lineHeight), color.White)
	//defaultFontInstance.Draw(g.imageBuf, resultText, int(0), int(rect.Min.Y))

	strs = strings.Split(resultText, "\n")
	y = 0
	for _, v := range strs {
		rect := defaultFontInstance.BoundString(v)
		y += rect.Dy() + margin
		defaultFontInstance.Draw(g.imageBuf, v, int(0), y)
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
