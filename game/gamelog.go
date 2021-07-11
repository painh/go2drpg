package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"image/color"
	"math"
)

type GameLogElement struct {
	text           string
	createdTs      int64
	key            string
	imageBuf       *ebiten.Image
	contentsHeight float64
	selectGroup    int
	selected       bool
	selectIndex    int
	drawRect       image.Rectangle
	mouseover      bool
	color          color.RGBA
}

//
//func (g *GameLogElement) Set(text string, rgba color.RGBA) {
//	g.selectGroup = 0
//	g.text = text
//	g.color = rgba
//
//	newText := ""
//	resultText := ""
//	linesCnt := 1
//	y := 0
//	for _, v := range text {
//		prevText := newText
//		newText += string(v)
//		rect := defaultFontInstance.BoundString(newText)
//		width := rect.Dx()
//		if width >= SettingConfigInstance.LogWidth {
//			rect := defaultFontInstance.BoundString(prevText)
//			y += rect.Dy() + SettingConfigInstance.LineSpacing
//
//			resultText = resultText + prevText + "\n"
//			prevText = newText
//			newText = string(v)
//			linesCnt++
//		}
//	}
//	resultText += newText
//	g.drawRect = defaultFontInstance.BoundString(newText)
//	y += g.drawRect.Dy() + SettingConfigInstance.LineSpacing
//	y += SettingConfigInstance.LineSpacing
//
//	strs := strings.Split(resultText, "\n")
//
//	g.contentsHeight = float64(y)
//	g.imageBuf = ebiten.NewImage(SettingConfigInstance.LogWidth, int(g.contentsHeight))
//
//	y = 0
//	for _, v := range strs {
//		rect := defaultFontInstance.BoundString(v)
//		y += rect.Dy() + SettingConfigInstance.LineSpacing
//		defaultFontInstance.Draw(g.imageBuf, v, 0, y, g.color)
//	}
//
//	DrawRect(g.imageBuf, 0, 0, float64(SettingConfigInstance.LogWidth), (g.contentsHeight), color.RGBA{255, 255, 255, 255})
//}
//

func (g *GameLogElement) Set(text string, rgba color.RGBA) {
	g.selectGroup = 0
	g.text = text
	g.color = rgba

	g.drawRect = defaultFontInstance.DrawWithWW(nil, text, 0, 0, g.color, true, float64(SettingConfigInstance.LogWidth))
	ay := math.Abs(float64(g.drawRect.Min.Y)) + float64(g.drawRect.Dy())
	g.contentsHeight = float64(ay) + float64(SettingConfigInstance.LineSpacing)
	g.imageBuf = ebiten.NewImage(SettingConfigInstance.LogWidth, int(g.contentsHeight))
	defaultFontInstance.DrawWithWW(g.imageBuf, text, 0, 0, g.color, false, float64(SettingConfigInstance.LogWidth))

	//DrawRect(g.imageBuf, 0, 0, float64(SettingConfigInstance.LogWidth), (g.contentsHeight), color.RGBA{255, 255, 255, 255})
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
		DrawRect(screen, x, y, float64(SettingConfigInstance.LogWidth), g.contentsHeight, color.RGBA{0, 255, 0, 255})
	}

	if g.selected {
		DrawRect(screen, x, y, float64(SettingConfigInstance.LogWidth), g.contentsHeight, color.White)
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

		if cursorY >= lineY-e.contentsHeight && cursorY < lineY {
			if e.selectGroup != g.currentSelectGroup {
				continue
			}

			if InputInstance.LBtnClicked() {
				g.LastSelectedIndex = e.selectIndex
				e.selected = true
				g.waitForSelect = false
				GameInstance.audio.PlayWave(SettingConfigInstance.WorkFolder + SettingConfigInstance.BtnClickSound)
				if g.callBack != nil {
					info := g.GetElementInfo(e.key)
					g.callBack(info)
				}
			} else {
				e.mouseover = true
			}
		}

		lineY -= float64(e.contentsHeight)
	}
}

func (g *GameLog) Draw(screen *ebiten.Image) {
	//g.logBuf.Clear()
	if GameInstance.status == GAME_UPDATE_STATUS_TALK_CHAR {
		g.logBuf.Fill(color.RGBA{0, 0, 0, 255})
	} else {
		g.logBuf.Fill(color.RGBA{64, 64, 64, 255})
	}

	y := float64(SettingConfigInstance.LogY + SettingConfigInstance.LogHeight)

	for i, e := range g.lines {
		if i >= SettingConfigInstance.LogLines {
			break
		}

		y -= float64(e.contentsHeight)
		e.Draw(g.logBuf, 0, y)
	}

	screen.DrawImage(g.logBuf, g.logBufOp)
}

func (g *GameLog) AddString(text string, rgba color.RGBA) {
	g.waitForSelect = false

	l := GameLogElement{text: text, createdTs: makeTimestamp()}
	l.Set(text, rgba)

	g.lines = append([]*GameLogElement{&l}, g.lines...)

	if len(g.lines) > SettingConfigInstance.LogLines {
		g.lines = g.lines[:len(g.lines)-1]
	}
}

func (g *GameLog) Add(rgba color.RGBA, a ...interface{}) {
	text := fmt.Sprint(a...)
	g.AddString(text, rgba)
}

func (g *GameLog) AddWithPrompt(rgba color.RGBA, a ...interface{}) {
	g.AddString(GameInstance.cursor.MakeTimePrompt(a...), rgba)
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
		l.Set(" • "+v.displayString, color.RGBA{0, 255, 0, 255})
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
	g.Add(color.RGBA{0, 255, 0, 255}, text)

	var list = []TextSelectElement{}

	list = append(list, TextSelectElement{key: "확인", displayString: "확인", info: "확인"})
	list = append(list, TextSelectElement{key: "취소", displayString: "취소", info: "취소"})
	GameInstance.log.TextSelect(list, func(info interface{}) {
		if info == "확인" {
			callBack()
		}
	})
}
