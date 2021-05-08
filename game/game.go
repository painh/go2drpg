package game

import (
	"fmt"
	"github.com/fardog/tmx"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/painh/go2drpg/assetmanager"
	"github.com/painh/go2drpg/game/scripts"
	"log"
	"path"
	"strconv"
)

var TILE_SIZE float64 = 32

const SPRITE_PATTERN = float64(16)

var SCALE float64 = TILE_SIZE / SPRITE_PATTERN

type Game struct {
	screenWidth            int
	screenHeight           int
	gameObjectManager      gameObjectManager
	FlowControllerInstance FlowController
	mapBuf                 *ebiten.Image
	mapBufOp               *ebiten.DrawImageOptions
	logBuf                 *ebiten.Image
	logBufOp               *ebiten.DrawImageOptions

	frameCnt     int64
	waitOneFrame int64
}

func (g *Game) WaitOneFrameOn() {
	g.waitOneFrame = g.frameCnt
}

func (g *Game) WaitOneFrame() {
	if g.waitOneFrame == 0 {
		return
	}

	if g.frameCnt == g.waitOneFrame {
		return
	}

	g.waitOneFrame = 0
	g.FlowControllerInstance.ShiftFlowToEventLoop()
}

func (g *Game) Update() error {
	g.frameCnt++
	InputInstance.Update()
	GameLogInstance.Update(InputInstance.x, InputInstance.y)

	g.WaitOneFrame()

	g.gameObjectManager.Update(InputInstance.x, InputInstance.y)

	dbClick := InputInstance.DBClick()

	if dbClick {
		fmt.Println(dbClick)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.mapBuf.Clear()
	g.logBuf.Clear()

	//g.mapBuf.Fill(color.RGBA{255, 0, 0, 255})
	//g.logBuf.Fill(color.RGBA{0, 255, 0, 255})

	g.gameObjectManager.Draw(g.mapBuf)

	//g.TextDialogInstance.Draw(screen)

	GameLogInstance.Draw(g.logBuf)

	screen.DrawImage(g.mapBuf, g.mapBufOp)
	screen.DrawImage(g.logBuf, g.logBufOp)

	fps := fmt.Sprintf("%f", ebiten.CurrentFPS())
	defaultFontInstance.DrawTextInBox(screen, "hello : "+fps, 0, 0)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

var GameInstance Game

func NewGame(screenWidth int, screenHeight int) *Game {
	ebiten.SetWindowSize(screenWidth, screenHeight)

	GameInstance = Game{}
	GameInstance.Init(screenWidth, screenHeight)

	scripts.Init(&GameInstance)

	assetmanager.Load("assets/16x16_Jerom_CC-BY-SA-3.0.png", "base")
	assetmanager.MakePatternImages("base", int(SPRITE_PATTERN), int(SPRITE_PATTERN))

	defaultFontInstance.LoadFont(ConfigInstance.FontPath, ConfigInstance.FontSize)

	file, err := ebitenutil.OpenFile((path.Join(".", "assets/tile.tmx")))
	if err != nil {
		log.Fatal(err.Error())
	}

	defer file.Close()

	m, err := tmx.Decode(file)
	if err != nil {
		log.Fatal(err.Error())
	}

	if m == nil {
		log.Fatal("map was nil")
	}

	if base := m.LayerWithName("base"); base != nil {
		trs, err := base.TileGlobalRefs()
		if err != nil {
			log.Fatal(err)
		}

		if l, e := len(trs), m.Width*m.Height; l != e {
			log.Fatalf("expected tiles of length %v, got %v", e, l)
		}

		tds, err := base.TileDefs(m.TileSets)
		if err != nil {
			log.Fatal(err)
		}

		cnt := 0

		for y := 0; y < base.Height; y++ {
			for x := 0; x < base.Width; x++ {
				GameInstance.gameObjectManager.GameSpriteAdd(float64(x), float64(y), TILE_SIZE, TILE_SIZE, "base_"+strconv.Itoa(int(tds[cnt].ID)))
				cnt++
			}
		}
	}

	if objects := m.ObjectGroupWithName("objLayer"); objects != nil {

		for _, v := range objects.Objects {
			x := v.X / SPRITE_PATTERN
			y := v.Y/SPRITE_PATTERN - 1

			GameInstance.gameObjectManager.GameObjectAdd(float64(x), float64(y), TILE_SIZE, TILE_SIZE, "base_"+strconv.Itoa(int(v.GlobalID-1)), v.Name)
		}

		//for y := 0; y < base.Height; y++ {
		//	for x := 0; x < base.Width; x++ {
		//		GameInstance.gameObjectManager.GameObjectAdd(float64(x*16), float64(y*16), "base_"+strconv.Itoa(int(tds[cnt].ID)))
		//		fmt.Println(tds[cnt])
		//		cnt++
		//	}
		//}
	}

	GameInstance.gameObjectManager.Width = float64(m.Width) * TILE_SIZE
	GameInstance.gameObjectManager.Height = float64(m.Height) * TILE_SIZE

	GameLogInstance.Add("클릭으로 선택, 더블클릭 혹은 우클릭으로 이동합니다.")
	GameLogInstance.Add("2")
	GameLogInstance.Add("3")
	GameLogInstance.Add("4")
	GameLogInstance.Add("5")

	return &GameInstance
}

func (g *Game) Init(screenWidth, screenHeight int) {
	g.screenHeight = screenHeight
	g.screenWidth = screenWidth

	g.gameObjectManager = gameObjectManager{}

	g.mapBuf = ebiten.NewImage(ConfigInstance.MapWidth, ConfigInstance.MapHeight)
	g.mapBufOp = &ebiten.DrawImageOptions{}
	g.mapBufOp.GeoM.Translate(float64(ConfigInstance.MapX), float64(ConfigInstance.MapY))

	g.logBuf = ebiten.NewImage(ConfigInstance.LogWidth, ConfigInstance.LogHeight)
	g.logBufOp = &ebiten.DrawImageOptions{}
	g.logBufOp.GeoM.Translate(float64(ConfigInstance.LogX), float64(ConfigInstance.LogY))

	g.FlowControllerInstance = FlowController{}
	g.FlowControllerInstance.Init()
}

func (g *Game) StartEvent() {
	g.FlowControllerInstance.StartEvent()
}

func (g *Game) ShiftFlowToMainLoop() {
	g.FlowControllerInstance.ShiftFlowToMainLoop()
}

func (g *Game) ShiftFlowToEventLoop() {
	g.FlowControllerInstance.ShiftFlowToEventLoop()
}

func (g *Game) EndEvent() {
	g.FlowControllerInstance.EventEnd()
}

func (g *Game) SetText(t string) {
	GameLogInstance.Add(t)
}

func (g *Game) TextSelect(t []string) {
	GameLogInstance.TextSelect(t)
}

func (g *Game) GetLastSelectedIndex() int {
	return GameLogInstance.LastSelectedIndex
}
