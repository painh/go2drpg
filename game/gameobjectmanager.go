package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/painh/go2drpg/assetmanager"
	"math"
)

type gameObjectManager struct {
	Width, Height          float64
	prevClickX, prevClickY float64
	tiles                  []*GameSprite
	objects                []*GameObject
}

func (g *gameObjectManager) Draw(screen *ebiten.Image) {
	for _, e := range g.tiles {
		e.Draw(screen)
	}

	for _, e := range g.tiles {
		e.Draw2(screen)
	}

	for _, e := range g.objects {
		e.Draw(screen)
	}

	for _, e := range g.objects {
		e.Draw2(screen)
	}
}

func (g *gameObjectManager) Refresh() {
	for _, e := range g.tiles {
		e.Refresh()
	}

	for _, e := range g.objects {
		e.Refresh()
	}
}

func (g *gameObjectManager) selectProcess(x, y float64) bool {
	objFound := false

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for _, e := range g.objects {
			e.selected = false
			if e.clickCheck(float64(x), float64(y)) {
				e.selected = true
				fmt.Println(e)
				objFound = true
			}
		}
	}

	return objFound
}

func (g *gameObjectManager) Update(x, y int) {
	objFound := false

	if DoubleClickInstance.LBtnPressed() {
		if !objFound {
			dx := g.prevClickX - float64(x)
			dy := g.prevClickY - float64(y)
			CameraInstance.SetXY(CameraInstance.x+dx, CameraInstance.y+dy)
			g.Refresh()
		}
	}

	if DoubleClickInstance.RBtnClicked() {
		worldX := float64(int(x+int(CameraInstance.x)) / int(TILE_SIZE))
		worldY := float64(int(y+int(CameraInstance.y)) / int(TILE_SIZE))

		for _, e := range g.objects {
			if e.selected {
				e.FindTo(worldX, worldY)
			}
		}

		//scripts.StartEvent("slime")
	}

	_, wy := ebiten.Wheel()
	if wy != 0 {
		TILE_SIZE = math.Max(1, TILE_SIZE+wy)
		SCALE = TILE_SIZE / SPRITE_PATTERN

		CameraInstance.Refresh()
		g.Refresh()
	}

	for _, e := range g.objects {
		e.Update()
	}

	g.prevClickX = float64(x)
	g.prevClickY = float64(y)
}

func (g *gameObjectManager) GameSpriteAdd(x, y, width, height float64, name string) {
	img := assetmanager.Get(name)
	obj := &GameSprite{x, y, 0, 0, width, height, false, img, &ebiten.DrawImageOptions{}}
	obj.Init()
	obj.SetXY(x, y)
	g.tiles = append(g.tiles, obj)
}

func (g *gameObjectManager) GameObjectAdd(x, y, width, height float64, sprName, objName string) {
	img := assetmanager.Get(sprName)
	obj := &GameObject{GameSprite: GameSprite{x, y, 0, 0, width, height, false, img, &ebiten.DrawImageOptions{}},
		cdmanager:   CooldownManager{dict: make(map[string]*Cooldown)},
		objName:     objName,
		movePosList: []*TilePos{}}
	//obj := &GameSprite{x, y, width, height, false, img, ebiten.DrawImageOptions{}}
	obj.SetXY(x, y)

	if objName == "Player" {
		obj.selected = true
	}

	obj.Init()
	g.objects = append(g.objects, obj)
}

func (g *gameObjectManager) CheckGameObjectPosition(x, y float64) bool {
	if x < 0 || y < 0 || x >= g.Width || y >= g.Height {
		return true
	}

	for _, e := range g.objects {
		if e.GameSprite.x == x && e.GameSprite.y == y {
			return true
		}
	}
	return false
}
