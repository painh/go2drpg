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
	"path"
	"strconv"
)

type Game struct {
	screenWidth            int
	screenHeight           int
	gameObjectManager      gameObjectManager
	TextDialogInstance     TextDialog
	FlowControllerInstance FlowController
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	DoubleClickInstance.Update(x, y, inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft))

	ok := g.FlowControllerInstance.WaitForEventLoop("")
	if !ok {
		g.gameObjectManager.Update(x, y)
		//fmt.Println("event loop가 죽었슴다")
	}

	g.TextDialogInstance.Update()

	dbClick := DoubleClickInstance.DBClick()

	if dbClick {
		fmt.Println(dbClick)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.gameObjectManager.Draw(screen)

	g.TextDialogInstance.Draw(screen)

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
	GameInstance.screenHeight = screenHeight
	GameInstance.screenWidth = screenWidth
	GameInstance.gameObjectManager = gameObjectManager{}
	GameInstance.TextDialogInstance = TextDialog{x: float64(screenWidth/2 - screenWidth/4),
		y: float64(screenHeight/2 - screenHeight/4)}
	GameInstance.TextDialogInstance.SetText(`
안녕
나의 이름은 김개똥
1
2
3
4
`)

	GameInstance.FlowControllerInstance = FlowController{}
	GameInstance.FlowControllerInstance.Init()

	scripts.Init(&GameInstance)

	assetmanager.Load("assets/16x16_Jerom_CC-BY-SA-3.0.png", "base")
	assetmanager.MakePatternImages("base", int(SPRITE_PATTERN), int(SPRITE_PATTERN))

	defaultFontInstance.LoadFont(ConfigInstance.Font_path, ConfigInstance.Font_size)

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
	return &GameInstance
}

func (g *Game) StartEvent() {
	g.FlowControllerInstance.StartEvent()
}

func (g *Game) WaitForMainLoop() {
	g.FlowControllerInstance.WaitForMainLoop()
}

func (g *Game) EndEvent() {
	g.FlowControllerInstance.EventEnd()
}

func (g *Game) SetText(t string) {
	g.TextDialogInstance.SetText(t)
}
