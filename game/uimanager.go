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
	color         color.RGBA
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
		GameInstance.audio.PlayWave(SettingConfigInstance.WorkFolder + SettingConfigInstance.BtnClickSound)
	}
}

func (u *UIButton) Draw(dst *ebiten.Image) {
	defaultFontInstance.DrawTextInBox(dst, u.text, u.x, u.y)
	DrawRect(dst, u.x, u.y, u.width, u.height, u.color)
}

type UIManager struct {
	uiDict  map[string]*UIWidget
	Clicked bool
}

func (u *UIManager) Init() {
	u.uiDict = map[string]*UIWidget{}

	u.AddButton("zoomout", float64(SettingConfigInstance.BtnZoomoutX), float64(SettingConfigInstance.BtnZoomoutY),
		float64(SettingConfigInstance.BtnZoomoutWidth), float64(SettingConfigInstance.BtnZoomoutHeight), "줌아웃", func() {
			SettingConfigInstance.RenderTileSize = math.Max(1, SettingConfigInstance.RenderTileSize+float64(SettingConfigInstance.ZoomStep))
			GameInstance.scale = SettingConfigInstance.RenderTileSize / SettingConfigInstance.RealTileSize
			GameInstance.cameraToCenter()
		}, color.RGBA{255, 255, 255, 255})

	u.AddButton("zoomin", float64(SettingConfigInstance.BtnZoominX), float64(SettingConfigInstance.BtnZoominY),
		float64(SettingConfigInstance.BtnZoominWidth), float64(SettingConfigInstance.BtnZoominHeight), "줌인", func() {
			SettingConfigInstance.RenderTileSize = math.Max(1, SettingConfigInstance.RenderTileSize-float64(SettingConfigInstance.ZoomStep))
			GameInstance.scale = SettingConfigInstance.RenderTileSize / SettingConfigInstance.RealTileSize
			GameInstance.cameraToCenter()
		}, color.RGBA{255, 255, 255, 255})

	u.AddButton("center", float64(SettingConfigInstance.BtnCenterX), float64(SettingConfigInstance.BtnCenterY),
		float64(SettingConfigInstance.BtnCenterWidth), float64(SettingConfigInstance.BtnCenterHeight), "중앙", func() {
			GameInstance.cameraToCenter()
		}, color.RGBA{255, 255, 255, 255})

	u.AddButton("person", float64(SettingConfigInstance.BtnPersonX), float64(SettingConfigInstance.BtnPersonY),
		float64(SettingConfigInstance.BtnPersonWidth), float64(SettingConfigInstance.BtnPersonHeight), "인물 목록", func() {
			GameInstance.log.AddWithPrompt("인물 목록")
		}, color.RGBA{255, 255, 255, 255})

	u.AddButton("place", float64(SettingConfigInstance.BtnLocationX), float64(SettingConfigInstance.BtnLocationY),
		float64(SettingConfigInstance.BtnLocationWidth), float64(SettingConfigInstance.BtnLocationHeight), "장소 목록", func() {
			GameInstance.log.Add("장소 목록 이동시간 : ", SettingConfigInstance.DefaultLocationMin)

			if len(GameInstance.player.activeLocation) == 0 {
				GameInstance.log.Add("활성 장소 없음")
			}

			var list = []TextSelectElement{}

			for _, v := range GameInstance.player.activeLocation {
				list = append(list, TextSelectElement{key: v.name, displayString: v.location.DisplayName, info: v})
			}
			GameInstance.log.TextSelect(list, func(info interface{}) {
				v := info.(Location)
				str := fmt.Sprintf("정말로 이동할까요? 이동에는 %v분이 소모됩니다.", SettingConfigInstance.DefaultLocationMin)
				GameInstance.log.Confirm(str, func() {
					GameInstance.log.AddWithPrompt("이동 ", v.location.DisplayName)
					GameInstance.player.AddTime(20)
					GameInstance.LoadMap(*v.location)
				})
			})
		}, color.RGBA{255, 255, 255, 255})
	u.AddButton("keyword", float64(SettingConfigInstance.BtnKeywordX), float64(SettingConfigInstance.BtnKeywordY),
		float64(SettingConfigInstance.BtnKeywordWidth), float64(SettingConfigInstance.BtnKeywordHeight), "키워드", func() {

			GameInstance.log.AddWithPrompt("키워드 목록")

			if len(GameInstance.keywordManager.dict) == 0 {
				GameInstance.log.Add("활성 키워드 없음")
			}

			var list = []TextSelectElement{}

			for k := range GameInstance.keywordManager.dict {
				list = append(list, TextSelectElement{key: k, displayString: k, info: k})
			}
			GameInstance.log.TextSelect(list, func(info interface{}) {
				GameInstance.player.AddTime(SettingConfigInstance.DefaultTalkMin)
				if GameInstance.scriptManager.RunKeywordScript(info.(string)) == false {
					GameInstance.log.AddWithPrompt(GameInstance.scriptManager.GetInvalidKeywordResponse())
				}
			})
		}, color.RGBA{255, 255, 255, 255})
	u.AddButton("talkend", float64(SettingConfigInstance.BtnTalkEndX), float64(SettingConfigInstance.BtnTalkEndY),
		float64(SettingConfigInstance.BtnTalkEndWidth), float64(SettingConfigInstance.BtnTalkEndHeight), "대화 종료", func() {
			GameInstance.TalkEnd()
		}, color.RGBA{255, 255, 255, 255})

	u.AddButton("devmode", 0, 20,
		100, 20, "스위치상태", func() {
			GameInstance.log.AddString("game switch")
			if len(GameInstance.gameSwitchManager.dict) == 0 {
				GameInstance.log.Add("활성 스위치 없음")
			}
			for k, v := range GameInstance.gameSwitchManager.dict {
				GameInstance.log.Add(k, " : ", v)
			}
		}, color.RGBA{0, 255, 0, 255})
}

func (u *UIManager) AddButton(name string, x, y, width, height float64, text string, onClick func(), color color.RGBA) {
	var btn UIWidget = &UIButton{x, y, width, height, text, onClick, color}
	u.uiDict[name] = &btn
}

func (u *UIManager) Update() {
	for _, e := range u.uiDict {
		(*e).Update()
	}
}

func (u *UIManager) Draw(dst *ebiten.Image) {
	for _, e := range u.uiDict {
		(*e).Draw(dst)
	}
}

func (u *UIManager) GetWidget(name string) *UIWidget {
	widget, exist := u.uiDict[name]

	if exist {
		return widget
	}

	return nil
}
