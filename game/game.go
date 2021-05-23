package game

import (
	"fmt"
	"github.com/fardog/tmx"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/painh/go2drpg/assetmanager"
	"github.com/painh/go2drpg/game/scripts"
	"log"
	"strconv"
)

var TILE_SIZE float64 = 32

var SPRITE_PATTERN = float64(16)

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
	//itemOriginManager      ItemOriginManager
	uimanager UIManager
	cursor    Cursor

	mapWidth  float64
	mapHeight float64

	frameCnt     int64
	waitOneFrame int64

	player Player
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
	g.uimanager.Clicked = false
	InputInstance.Update()
	g.uimanager.Update()
	g.cursor.Update()
	GameLogInstance.Update(InputInstance.x, InputInstance.y)

	g.WaitOneFrame()

	if !g.uimanager.Clicked {
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
		x = x / float64(len(list)) * TILE_SIZE
		y = y / float64(len(list)) * TILE_SIZE
		CameraInstance.SetXY(x-float64(ConfigInstance.MapWidth/2)+TILE_SIZE/2, y-float64(ConfigInstance.MapHeight/2)+TILE_SIZE/2)

		g.gameObjectManager.Refresh()
	}
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

	GameInstance = Game{}
	GameInstance.Init(screenWidth, screenHeight)

	scripts.Init(&GameInstance)

	defaultFontInstance.LoadFont(ConfigInstance.FontPath, ConfigInstance.FontSize)

	//GameInstance.itemOriginManager = ItemOriginManager{dict: make(map[int]*ItemOrigin)}
	//GameInstance.itemOriginManager.LoadFromCSV("assets/items.csv")

	SPRITE_PATTERN = float64(ConfigInstance.SpritePatternSize)
	assetmanager.Load(ConfigInstance.TileSpriteFilename, "base")
	assetmanager.MakePatternImages("base", int(SPRITE_PATTERN), int(SPRITE_PATTERN))

	GameInstance.LoadMap(ConfigInstance.LocationList[0])

	GameLogInstance.Add("클릭으로 선택, 더블클릭 혹은 우클릭으로 이동합니다.")

	GameInstance.uimanager.Init()
	GameInstance.cursor.Init()
	GameInstance.player.Init()
	GameInstance.player.ActiveLocation(ConfigInstance.LocationList[0].Name)
	GameInstance.player.ActiveLocation(ConfigInstance.LocationList[1].Name)

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
	log.Fatal("현재는 쓰지 않는 기능")
	//GameLogInstance.TextSelect(t)
}

func (g *Game) GetLastSelectedIndex() int {
	return GameLogInstance.LastSelectedIndex
}

func (g *Game) LoadMap(info LocationInfo) {
	file, err := ebitenutil.OpenFile(info.Filename)
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

	g.gameObjectManager.Clear()

	if base := m.LayerWithName("base"); base != nil {
		trs, err := base.TileGlobalRefs()
		if err != nil {
			log.Fatal(err)
		}

		g.mapWidth = float64(m.Width)
		g.mapHeight = float64(m.Height)

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
				g.gameObjectManager.GameSpriteAdd(float64(x), float64(y), SPRITE_PATTERN, SPRITE_PATTERN, "base_"+strconv.Itoa(int(tds[cnt].ID)))
				cnt++
			}
		}
	}

	if objects := m.ObjectGroupWithName("obj"); objects != nil {

		for _, v := range objects.Objects {
			x := v.X / SPRITE_PATTERN
			y := v.Y / SPRITE_PATTERN

			if v.GlobalID != 0 { //GlobalID가 없다면, 이미지가 선택되지 않은 충돌 사각형으로 생각하여 원본을 씀. Tiled가 오브젝트일때는 y좌표를 + 1해서 주는 방식이라 후처리를 해야함.
				y--
			}

			g.gameObjectManager.GameObjectAdd(float64(x), float64(y), v.Width, v.Height, "base_"+strconv.Itoa(int(v.GlobalID-1)), v.Name)
		}
	}

	g.gameObjectManager.Width = float64(m.Width) * TILE_SIZE
	g.gameObjectManager.Height = float64(m.Height) * TILE_SIZE
	g.cameraToCenter()
}
