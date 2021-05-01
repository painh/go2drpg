package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/painh/go2drpg/assetmanager"
	"image/color"
)

type GameSprite struct {
	x        float64
	y        float64
	screenX  float64
	screenY  float64
	width    float64
	height   float64
	selected bool

	img *assetmanager.ImageResource
	op  *ebiten.DrawImageOptions
}

func (g *GameSprite) Init() {
	g.op.GeoM.Reset()
	g.SetXY(g.x, g.y)

}

func (g *GameSprite) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.img.Img, g.op)

	//if g.selected {
	//	DrawRect(screen, g.x, g.y, g.width, g.height, color.RGBA{0, 255, 0, 255})
	//}
}

func (g *GameSprite) Draw2(screen *ebiten.Image) {
	//screen.DrawImage(g.img.Img, &g.op)

	if g.selected {
		DrawRect(screen, g.screenX, g.screenY, TILE_SIZE, TILE_SIZE, color.RGBA{G: 255, A: 255})
	}
}

func (g *GameSprite) clickCheck(x, y float64) bool {
	if x >= g.screenX && y >= g.screenY && g.screenX+g.width > x && g.screenY+g.height > y {
		return true
	}

	return false
}

func (g *GameSprite) Refresh() {
	g.op = &ebiten.DrawImageOptions{}
	g.screenX = (g.x * TILE_SIZE) - CameraInstance.x
	g.screenY = (g.y * TILE_SIZE) - CameraInstance.y

	g.op.GeoM.Translate(g.screenX/SCALE, g.screenY/SCALE)
	g.op.GeoM.Scale(SCALE, SCALE)
	g.width = TILE_SIZE
	g.height = TILE_SIZE
}

func (g *GameSprite) SetXY(x, y float64) {
	g.x = x
	g.y = y
	g.Refresh()
}
