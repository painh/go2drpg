package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/painh/go2drpg/assetmanager"
	"image/color"
)

type GameSprite struct {
	x      float64
	y      float64
	width  float64
	height float64

	screenX      float64
	screenY      float64
	screenWidth  float64
	screenHeight float64
	selected     bool

	img *assetmanager.ImageResource
	op  *ebiten.DrawImageOptions
}

func (g *GameSprite) Init() {
	g.op.GeoM.Reset()
	g.SetXY(g.x, g.y)

}

func (g *GameSprite) Draw(screen *ebiten.Image) {
	if g.img != nil {
		screen.DrawImage(g.img.Img, g.op)
	} else {
		DrawRect(screen, g.screenX, g.screenY, g.screenWidth, g.screenHeight, color.RGBA{R: 255, A: 255})
	}
}

func (g *GameSprite) clickCheck(x, y float64) bool {
	if x >= g.screenX && y >= g.screenY && g.screenX+g.screenWidth > x && g.screenY+g.screenHeight > y {
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
	g.screenWidth = g.width * SCALE
	g.screenHeight = g.height * SCALE
}

func (g *GameSprite) SetXY(x, y float64) {
	g.x = x
	g.y = y
	g.Refresh()
}
