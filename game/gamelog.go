package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type GameLogElement struct {
	text      string
	createdTs int64
	rect      image.Rectangle
}

func (g *GameLogElement) Set(text string) {
	g.text = text
	g.rect = defaultFontInstance.BoundString(text)
}

func (g *GameLogElement) Update() {

}

func (g *GameLogElement) Draw(screen *ebiten.Image, x, y float64) {
	defaultFontInstance.DrawTextInBox(screen, g.text, x, y)

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
		if i >= ConfigInstance.Log_lines {
			break
		}

		y -= float64(e.rect.Dy())
		e.Draw(screen, 0, y)
	}
}

func (g *GameLog) Add(text string) {
	l := GameLogElement{text: text, createdTs: makeTimestamp()}
	l.Set(text)

	g.lines = append([]GameLogElement{l}, g.lines...)

	if len(g.lines) > ConfigInstance.Log_lines {
		g.lines = g.lines[:len(g.lines)-1]
	}
}
