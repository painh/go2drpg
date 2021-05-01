package game

import (
	"fmt"
	"github.com/beefsack/go-astar"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type GameObject struct {
	GameSprite
	cdmanager CooldownManager
	objName   string

	movePosList []*TilePos
}

func (g *GameObject) Init() {
	g.GameSprite.Init()
}

func (g *GameObject) FindTo(x, y float64) {
	from := tileManager.Get(g.x, g.y)
	to := tileManager.Get(x, y)
	path, distance, found := astar.Path(from, to)
	fmt.Println(path)
	fmt.Println(distance)
	fmt.Println(found)

	if !found {
		return
	}

	g.movePosList = []*TilePos{}
	for i := len(path) - 1; i >= 0; i-- {
		v := path[i]
		g.movePosList = append(g.movePosList, v.(*TilePos))
	}
}

func (g *GameObject) Draw2(screen *ebiten.Image) {
	g.GameSprite.Draw2(screen)

	for _, v := range g.movePosList {
		DrawRect(screen, v.x*TILE_SIZE-CameraInstance.x, v.y*TILE_SIZE-CameraInstance.y, TILE_SIZE, TILE_SIZE, color.RGBA{0, 255, 0, 255})
	}

}

func (g *GameObject) Update() {
	if len(g.movePosList) == 0 {
		return
	}

	if g.cdmanager.IsCooldownOver("move", 500) {
		g.cdmanager.ActiveCooldown("move")
		tile := g.movePosList[0]
		g.movePosList = g.movePosList[1:]
		g.GameSprite.SetXY(tile.x, tile.y)
	}
}
