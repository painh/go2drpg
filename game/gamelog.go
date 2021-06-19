package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"image/color"
	"strings"
)

type GameLogElement struct {
	text        string
	createdTs   int64
	key         string
	imageBuf    *ebiten.Image
	lineHeight  float64
	selectGroup int
	selected    bool
	selectIndex int
	drawRect    image.Rectangle
	mouseover   bool
}

func (g *GameLogElement) Set(text string) {
	g.selectGroup = 0
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
		if width >= SettingConfigInstance.LogWidth {
			rect := defaultFontInstance.BoundString(prevText)
			y += rect.Dy() + SettingConfigInstance.LineSpacing

			resultText = prevText + "\n"
			prevText = newText
			newText = string(v)
			linesCnt++
		}
	}
	resultText += newText
	g.drawRect = defaultFontInstance.BoundString(newText)
	y += g.drawRect.Dy() + SettingConfigInstance.LineSpacing
	y += SettingConfigInstance.LineSpacing

	strs := strings.Split(resultText, "\n")

	g.lineHeight = float64(y)
	g.imageBuf = ebiten.NewImage(SettingConfigInstance.LogWidth, int(g.lineHeight))

	y = 0
	for _, v := range strs {
		rect := defaultFontInstance.BoundString(v)
		y += rect.Dy() + SettingConfigInstance.LineSpacing
		defaultFontInstance.Draw(g.imageBuf, v, 0, y)
	}

}

func (g *GameLogElement) Update() {

}

func (g *GameLogElement) Draw(screen *ebiten.Image, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	if g.selectGroup != 0 {
		ebitenutil.DrawRect(screen, x, y+1, float64(SettingConfigInstance.LogWidth), float64(g.drawRect.Dy()+SettingConfigInstance.LineSpacing*2)-2, color.RGBA{128, 128, 128, 255})
	}
	screen.DrawImage(g.imageBuf, op)
	if g.mouseover {
		DrawRect(screen, x, y, float64(SettingConfigInstance.LogWidth), g.lineHeight, color.RGBA{0, 255, 0, 255})
	}

	if g.selected {
		DrawRect(screen, x, y, float64(SettingConfigInstance.LogWidth), g.lineHeight, color.White)
	}
}

type GameLog struct {
	lines              []*GameLogElement
	currentSelectGroup int
	waitForSelect      bool
	LastSelectedIndex  int
	callBack           func(info interface{})
	textSelectElement  []TextSelectElement

	logBuf   *ebiten.Image
	logBufOp *ebiten.DrawImageOptions
}

func (g *GameLog) GetElementInfo(key string) interface{} {
	for _, v := range g.textSelectElement {
		if v.key == key {
			return v.info
		}
	}

	return nil
}

func (g *GameLog) Update(x, y int) {
	if !g.waitForSelect {
		return
	}
	//
	//if !(x >= SettingConfigInstance.LogX && x <= SettingConfigInstance.LogX+SettingConfigInstance.LogWidth &&
	//	y >= SettingConfigInstance.LogY && y <= SettingConfigInstance.LogY+SettingConfigInstance.LogHeight) {
	//	return
	//}

	//cursorX := float64(x - SettingConfigInstance.LogX)
	cursorY := float64(y - SettingConfigInstance.LogY)

	lineY := float64(SettingConfigInstance.LogY + SettingConfigInstance.LogHeight)

	for i, e := range g.lines {
		if i >= SettingConfigInstance.LogLines {
			break
		}

		e.mouseover = false

		if cursorY >= lineY-e.lineHeight && cursorY < lineY {
			if e.selectGroup != g.currentSelectGroup {
				continue
			}

			if InputInstance.LBtnClicked() {
				g.LastSelectedIndex = e.selectIndex
				e.selected = true
				g.waitForSelect = false
				if g.callBack != nil {
					info := g.GetElementInfo(e.key)
					g.callBack(info)
				}
			} else {
				e.mouseover = true
			}
		}

		lineY -= float64(e.lineHeight)
	}
}

func (g *GameLog) Draw(screen *ebiten.Image) {
	g.logBuf.Clear()

	y := float64(SettingConfigInstance.LogY + SettingConfigInstance.LogHeight)

	for i, e := range g.lines {
		if i >= SettingConfigInstance.LogLines {
			break
		}

		y -= float64(e.lineHeight)
		e.Draw(g.logBuf, 0, y)
	}

	screen.DrawImage(g.logBuf, g.logBufOp)
}

func (g *GameLog) AddString(text string) {
	g.waitForSelect = false

	l := GameLogElement{text: text, createdTs: makeTimestamp()}
	l.Set(text)

	g.lines = append([]*GameLogElement{&l}, g.lines...)

	if len(g.lines) > SettingConfigInstance.LogLines {
		g.lines = g.lines[:len(g.lines)-1]
	}
}

func (g *GameLog) Add(a ...interface{}) {
	text := fmt.Sprint(a...)
	g.AddString(text)
}

func (g *GameLog) AddWithPrompt(a ...interface{}) {
	g.AddString(GameInstance.cursor.MakeTimePrompt(a...))
}

type TextSelectElement struct {
	displayString string
	key           string
	info          interface{}
}

func (g *GameLog) TextSelect(t []TextSelectElement, callBack func(info interface{})) {

	g.currentSelectGroup++
	g.waitForSelect = true
	g.callBack = callBack

	for i, v := range t {
		l := GameLogElement{text: "", createdTs: makeTimestamp(), key: v.key}
		l.Set(" • " + v.displayString)
		l.selectGroup = g.currentSelectGroup
		l.selectIndex = i

		g.lines = append([]*GameLogElement{&l}, g.lines...)

		if len(g.lines) > SettingConfigInstance.LogLines {
			g.lines = g.lines[:len(g.lines)-1]
		}
	}

	g.textSelectElement = t
}

func (g *GameLog) Confirm(text string, callBack func()) {
	g.Add(text)

	var list = []TextSelectElement{}

	list = append(list, TextSelectElement{key: "확인", displayString: "확인", info: "확인"})
	list = append(list, TextSelectElement{key: "취소", displayString: "취소", info: "취소"})
	GameInstance.log.TextSelect(list, func(info interface{}) {
		if info == "확인" {
			callBack()
		}
	})
}
