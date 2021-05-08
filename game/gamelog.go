package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"strings"
)

type GameLogElement struct {
	text        string
	createdTs   int64
	imageBuf    *ebiten.Image
	lineHeight  float64
	selectGroup int
	selected    bool
	selectIndex int
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

	g.lineHeight = float64(y)
	g.imageBuf = ebiten.NewImage(ConfigInstance.LogWidth, int(g.lineHeight))

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

	if g.selected {
		DrawRect(screen, x, y, float64(ConfigInstance.LogWidth), g.lineHeight, color.White)
	}
}

type GameLog struct {
	lines              []*GameLogElement
	currentSelectGroup int
	waitForSelect      bool
	LastSelectedIndex  int
}

var GameLogInstance = GameLog{lines: []*GameLogElement{}}

func (g *GameLog) Update(x, y int) {
	if !g.waitForSelect {
		return
	}

	if !InputInstance.LBtnClicked() {
		return
	}

	if !(x >= ConfigInstance.LogX && x <= ConfigInstance.LogX+ConfigInstance.LogWidth &&
		y >= ConfigInstance.LogY && y <= ConfigInstance.LogY+ConfigInstance.LogHeight) {
		return
	}

	//cursorX := float64(x - ConfigInstance.LogX)
	cursorY := float64(y - ConfigInstance.LogY)

	lineY := float64(ConfigInstance.LogY + ConfigInstance.LogHeight)

	for i, e := range g.lines {
		if i >= ConfigInstance.LogLines {
			break
		}

		if cursorY >= lineY-e.lineHeight && cursorY < lineY {
			if e.selectGroup != g.currentSelectGroup {
				break
			}
			g.LastSelectedIndex = e.selectIndex
			e.selected = true
			g.waitForSelect = false
			GameInstance.ShiftFlowToEventLoop()
		}

		lineY -= float64(e.lineHeight)
	}
}

func (g *GameLog) Draw(screen *ebiten.Image) {
	y := float64(ConfigInstance.LogY + ConfigInstance.LogHeight)

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

	g.lines = append([]*GameLogElement{&l}, g.lines...)

	if len(g.lines) > ConfigInstance.LogLines {
		g.lines = g.lines[:len(g.lines)-1]
	}

	fmt.Println(text)
}

func (g *GameLog) TextSelect(t []string) {

	g.currentSelectGroup++
	g.waitForSelect = true

	for i, v := range t {
		l := GameLogElement{text: "", createdTs: makeTimestamp()}
		l.Set(" • " + v)
		l.selectGroup = g.currentSelectGroup
		l.selectIndex = i

		g.lines = append([]*GameLogElement{&l}, g.lines...)

		if len(g.lines) > ConfigInstance.LogLines {
			g.lines = g.lines[:len(g.lines)-1]
		}
	}
}
