package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type UIWidget interface {
	Update()
	Draw(dst *ebiten.Image)
}

type UIButton struct {
	x, y          float64
	width, height float64
	text          string
	OnClick       func()
}

func (u *UIButton) SetText(text string) {
	u.text = text
}

func (u *UIButton) Update() {
	if !InputInstance.LBtnClicked() {
		return
	}

	if InputInstance.x >= int(u.x) && InputInstance.y >= int(u.y) &&
		InputInstance.x < int(u.x+u.width) && InputInstance.y < int(u.y+u.height) {
		u.OnClick()
	}
}

func (u *UIButton) Draw(dst *ebiten.Image) {
	defaultFontInstance.DrawTextInBox(dst, u.text, u.x, u.y)
	DrawRect(dst, u.x, u.y, u.width, u.height, color.RGBA{255, 255, 255, 255})
}

type UIManager struct {
	uilist []UIWidget
}

func (u *UIManager) Init() {
	u.uilist = []UIWidget{}

	u.AddButton(float64(ConfigInstance.BtnPersonX), float64(ConfigInstance.BtnPersonY),
		float64(ConfigInstance.BtnPersonWidth), float64(ConfigInstance.BtnPersonHeight), "인물 목록", func() {
			fmt.Println("사람목록")
		})

	u.AddButton(float64(ConfigInstance.BtnLocationX), float64(ConfigInstance.BtnLocationY),
		float64(ConfigInstance.BtnLocationWidth), float64(ConfigInstance.BtnLocationHeight), "장소 목록", func() {
			fmt.Println("장소 목록")
		})
	u.AddButton(float64(ConfigInstance.BtnItemX), float64(ConfigInstance.BtnItemY),
		float64(ConfigInstance.BtnItemWidth), float64(ConfigInstance.BtnItemHeight), "아이템 목록", func() {
			fmt.Println("아이템 목록")
		})
}

func (u *UIManager) AddButton(x, y, width, height float64, text string, onClick func()) {
	var btn UIWidget = &UIButton{x, y, width, height, text, onClick}
	u.uilist = append(u.uilist, btn)
}

func (u *UIManager) Update() {
	for _, e := range u.uilist {
		e.Update()
	}
}

func (u *UIManager) Draw(dst *ebiten.Image) {
	for _, e := range u.uilist {
		e.Draw(dst)
	}
}
