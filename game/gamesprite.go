package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/painh/go2drpg/assetmanager"
)

type GameSprite struct {
	x       float64
	y       float64
	width   float64
	height  float64
	offsetY float64

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
	}
	//else {
	//	DrawRect(screen, g.screenX, g.screenY, g.screenWidth, g.screenHeight, color.RGBA{R: 255, A: 255})
	//}
}

func (g *GameSprite) clickCheck(x, y float64) bool {
	if x >= g.screenX && y >= g.screenY && g.screenX+g.screenWidth > x && g.screenY+g.screenHeight > y {
		return true
	}

	return false
}

func (g *GameSprite) Refresh() {
	g.op = &ebiten.DrawImageOptions{}
	g.screenX = (g.x * SettingConfigInstance.RenderTileSize) - CameraInstance.x
	g.screenY = ((g.y + g.offsetY) * SettingConfigInstance.RenderTileSize) - CameraInstance.y

	g.op.GeoM.Translate(g.screenX/GameInstance.scale, g.screenY/GameInstance.scale)
	g.op.GeoM.Scale(GameInstance.scale, GameInstance.scale)
	g.screenWidth = g.width * GameInstance.scale
	g.screenHeight = g.height * GameInstance.scale
}

func (g *GameSprite) SetXY(x, y float64) {
	g.x = x
	g.y = y
	g.Refresh()
}
