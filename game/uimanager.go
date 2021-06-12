package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
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
		GameInstance.uimanager.Clicked = true
	}
}

func (u *UIButton) Draw(dst *ebiten.Image) {
	defaultFontInstance.DrawTextInBox(dst, u.text, u.x, u.y)
	DrawRect(dst, u.x, u.y, u.width, u.height, color.RGBA{255, 255, 255, 255})
}

type UIManager struct {
	uilist  []UIWidget
	Clicked bool
}

func (u *UIManager) Init() {
	u.uilist = []UIWidget{}

	u.AddButton(float64(SettingConfigInstance.BtnZoomoutX), float64(SettingConfigInstance.BtnZoomoutY),
		float64(SettingConfigInstance.BtnZoomoutWidth), float64(SettingConfigInstance.BtnZoomoutHeight), "줌아웃", func() {
			SettingConfigInstance.RenderTileSize = math.Max(1, SettingConfigInstance.RenderTileSize+float64(SettingConfigInstance.ZoomStep))
			GameInstance.scale = SettingConfigInstance.RenderTileSize / SettingConfigInstance.RealTileSize
			GameInstance.cameraToCenter()
		})

	u.AddButton(float64(SettingConfigInstance.BtnZoominX), float64(SettingConfigInstance.BtnZoominY),
		float64(SettingConfigInstance.BtnZoominWidth), float64(SettingConfigInstance.BtnZoominHeight), "줌인", func() {
			SettingConfigInstance.RenderTileSize = math.Max(1, SettingConfigInstance.RenderTileSize-float64(SettingConfigInstance.ZoomStep))
			GameInstance.scale = SettingConfigInstance.RenderTileSize / SettingConfigInstance.RealTileSize
			GameInstance.cameraToCenter()
		})

	u.AddButton(float64(SettingConfigInstance.BtnCenterX), float64(SettingConfigInstance.BtnCenterY),
		float64(SettingConfigInstance.BtnCenterWidth), float64(SettingConfigInstance.BtnCenterHeight), "중앙", func() {
			GameInstance.cameraToCenter()
		})

	u.AddButton(float64(SettingConfigInstance.BtnPersonX), float64(SettingConfigInstance.BtnPersonY),
		float64(SettingConfigInstance.BtnPersonWidth), float64(SettingConfigInstance.BtnPersonHeight), "인물 목록", func() {
			fmt.Println("사람목록")
		})

	u.AddButton(float64(SettingConfigInstance.BtnLocationX), float64(SettingConfigInstance.BtnLocationY),
		float64(SettingConfigInstance.BtnLocationWidth), float64(SettingConfigInstance.BtnLocationHeight), "장소 목록", func() {
			var list []TextSelectElement = []TextSelectElement{}

			for _, v := range GameInstance.player.activeLocation {
				list = append(list, TextSelectElement{key: v.name, displayString: v.location.DisplayName, info: v})
			}
			GameInstance.Log.TextSelect(list, func(info interface{}) {
				v := info.(Location)
				GameInstance.Log.AddWithPrompt("이동 ", v.location.DisplayName)
				GameInstance.player.AddTime(20)
				GameInstance.LoadMap(*v.location)
			})
		})
	u.AddButton(float64(SettingConfigInstance.BtnItemX), float64(SettingConfigInstance.BtnItemY),
		float64(SettingConfigInstance.BtnItemWidth), float64(SettingConfigInstance.BtnItemHeight), "아이템 목록", func() {
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
