package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
)

const GAME_UPDATE_STATUS_MAP_INTERACTION = 0
const GAME_UPDATE_STATUS_WAIT_USER_LOG_INTERACTION = 1
const GAME_UPDATE_STATUS_TALK_CHAR = 2

type Game struct {
	status            int
	screenWidth       int
	screenHeight      int
	gameObjectManager gameObjectManager
	mapBuf            *ebiten.Image
	mapBufOp          *ebiten.DrawImageOptions
	log               *GameLog
	//itemOriginManager      ItemOriginManager
	uimanager UIManager
	cursor    Cursor

	mapWidth  float64
	mapHeight float64

	frameCnt     int64
	waitOneFrame int64

	player Player

	music MusicManager

	scale float64

	mapLoader MapLoader

	scriptManager     ScriptManager
	keywordManager    KeywordManager
	gameSwitchManager GameSwitchManager
}

func (g *Game) WaitOneFrameOn() {
	g.waitOneFrame = g.frameCnt
}

func (g *Game) Update() error {
	g.frameCnt++
	g.uimanager.Clicked = false
	InputInstance.Update()
	g.uimanager.Update()
	g.cursor.Update()
	g.scriptManager.Update()
	g.log.Update(InputInstance.x, InputInstance.y)
	g.music.Update()

	if InputInstance.RBtnClicked() && g.status != GAME_UPDATE_STATUS_MAP_INTERACTION && g.gameObjectManager.Inbound(InputInstance.x, InputInstance.y) {
		g.log.AddWithPrompt("대화를 종료해 주세요.")
	}

	if !g.uimanager.Clicked && g.status == GAME_UPDATE_STATUS_MAP_INTERACTION {
		g.gameObjectManager.Update(InputInstance.x, InputInstance.y)
	}

	//dbClick := InputInstance.DBClick()
	//
	//if dbClick {
	//	fmt.Println(dbClick)
	//}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.cameraToCenter()
	}

	return nil
}

func (g *Game) cameraToCenter() {
	list := g.gameObjectManager.GetSelectedList()
	if len(list) > 0 {
		x := float64(0)
		y := float64(0)
		for i := 0; i < len(list); i++ {
			x += list[i].x
			y += list[i].y
		}
		x = x / float64(len(list)) * SettingConfigInstance.RenderTileSize
		y = y / float64(len(list)) * SettingConfigInstance.RenderTileSize
		CameraInstance.SetXY(x-float64(SettingConfigInstance.MapWidth/2)+SettingConfigInstance.RenderTileSize/2, y-float64(SettingConfigInstance.MapHeight/2)+SettingConfigInstance.RenderTileSize/2)

		g.gameObjectManager.Refresh()
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.mapBuf.Clear()

	//g.mapBuf.Fill(color.RGBA{255, 0, 0, 255})
	//g.logBuf.Fill(color.RGBA{0, 255, 0, 255})

	g.gameObjectManager.Draw(g.mapBuf)

	//g.TextDialogInstance.Draw(screen)

	g.log.Draw(screen)

	screen.DrawImage(g.mapBuf, g.mapBufOp)
	g.uimanager.Draw(screen)
	g.cursor.Draw(screen)

	fps := fmt.Sprintf("%f", ebiten.CurrentFPS())
	defaultFontInstance.DrawTextInBox(screen, "hello : "+fps, 0, 0)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

var GameInstance Game

func NewGame(screenWidth int, screenHeight int) *Game {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	defaultFontInstance.LoadFont(SettingConfigInstance.FontPath, SettingConfigInstance.FontSize)

	GameInstance = Game{}
	GameInstance.Init(screenWidth, screenHeight)

	return &GameInstance
}

func (g *Game) Init(screenWidth, screenHeight int) {
	g.screenHeight = screenHeight
	g.screenWidth = screenWidth

	g.gameObjectManager = gameObjectManager{}

	g.mapBuf = ebiten.NewImage(SettingConfigInstance.MapWidth, SettingConfigInstance.MapHeight)
	g.mapBufOp = &ebiten.DrawImageOptions{}
	g.mapBufOp.GeoM.Translate(float64(SettingConfigInstance.MapX), float64(SettingConfigInstance.MapY))

	g.scale = SettingConfigInstance.RenderTileSize / SettingConfigInstance.RealTileSize
	g.LoadMap(SettingConfigInstance.LocationList[0])
	g.log = &GameLog{lines: []*GameLogElement{}}
	g.log.logBuf = ebiten.NewImage(SettingConfigInstance.LogWidth, SettingConfigInstance.LogHeight)
	g.log.logBufOp = &ebiten.DrawImageOptions{}
	g.log.logBufOp.GeoM.Translate(float64(SettingConfigInstance.LogX), float64(SettingConfigInstance.LogY))

	g.log.Add("클릭으로 선택, 더블클릭 혹은 우클릭으로 이동합니다.")

	g.uimanager.Init()
	g.cursor.Init()
	g.player.Init()
	g.player.ActiveLocation(SettingConfigInstance.LocationList[0].Name)
	g.player.ActiveLocation(SettingConfigInstance.LocationList[1].Name)

	g.status = GAME_UPDATE_STATUS_MAP_INTERACTION
	g.scriptManager.Init()
	g.keywordManager.Init()
	g.gameSwitchManager.Init()

	for _, v := range SettingConfigInstance.Scripts {
		ScriptLoad(SettingConfigInstance.WorkFolder + v)
	}
}

func (g *Game) SetText(t string) {
	g.log.Add(t)
}

func (g *Game) TextSelect(t []string) {
	log.Fatal("현재는 쓰지 않는 기능")
	//GameLogInstance.TextSelect(t)
}

func (g *Game) GetLastSelectedIndex() int {
	return g.log.LastSelectedIndex
}

func (g *Game) LoadMap(info LocationInfo) {
	g.mapLoader.Load(info.Filename, g)

	g.gameObjectManager.Width = float64(g.mapLoader.Width) * SettingConfigInstance.RenderTileSize
	g.gameObjectManager.Height = float64(g.mapLoader.Height) * SettingConfigInstance.RenderTileSize
	g.cameraToCenter()

	g.music.Init()
}

func (g *Game) SetStatus(status int) {
	g.status = status
	//switch status {
	//case 0:
	//	g.status = status
	//case 1:
	//	g.status = status
	//case 2:
	//	g.status = status
	//
	//}
}
func (g *Game) TalkEnd() {
	if g.status == GAME_UPDATE_STATUS_MAP_INTERACTION {
		g.log.AddWithPrompt("지금은 대화중이 아닙니다.")
		return
	}

	g.SetStatus(GAME_UPDATE_STATUS_MAP_INTERACTION)
	g.log.AddWithPrompt("대화를 종료합니다.")
	g.gameObjectManager.SetActiveObject(nil)
}
