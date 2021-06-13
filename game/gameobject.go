package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/painh/go2drpg/game/scripts"
	"image/color"
)

type GameObject struct {
	GameSprite
	cdmanager CooldownManager
	objName   string
	zindex    float64

	movePosList []*TilePos
}

func (g *GameObject) Init() {
	g.GameSprite.Init()
}

func (g *GameObject) SetSize(width, height float64) {
	g.width = width
	g.height = height
}

func (g *GameObject) FindTo(x, y float64) {
	path, _, found := tileManagerInstance.FindTo(g.x, g.y, x, y, g)

	if !found {
		GameInstance.Log.AddWithPrompt("can't move there")
		return
	}

	g.movePosList = []*TilePos{}
	for i := len(path) - 2; i >= 0; i-- { // 첫위치는 제외함. -2.
		v := path[i]
		g.movePosList = append(g.movePosList, v.(*TilePos))
	}

	GameInstance.Log.AddWithPrompt("Move : ", x, y)

}

func (g *GameObject) Update() {
	if len(g.movePosList) == 0 {
		return
	}

	if g.cdmanager.IsCooldownOver("move", 500) {
		g.cdmanager.ActiveCooldown("move")
		tile := g.movePosList[0]
		obj := GameInstance.gameObjectManager.CheckGameObjectPosition(tile.x*SettingConfigInstance.RealTileSize, tile.y*SettingConfigInstance.RealTileSize, g.zindex, g.width, g.height, g)
		if obj != nil {
			if len(g.movePosList) == 1 {
				//도착점이면 이벤트 발생
				g.movePosList = g.movePosList[1:]

				if GameInstance.scriptManager.RunObjectScript(obj.objName) {
					//ok
				} else {
					GameInstance.Log.AddWithPrompt("아무것도 없습니다.")
				}
				return
			}
		}

		g.movePosList = g.movePosList[1:]
		g.GameSprite.SetXY(tile.x, tile.y)

		if g.objName == SettingConfigInstance.PlayerObjectName {
			GameInstance.player.AddTime(SettingConfigInstance.DefaultMoveMin)
		}
	}
}

func (g *GameObject) CheckCollision(x, y, width, height float64) bool {
	if g.x*SettingConfigInstance.RealTileSize < x+width &&
		g.x*SettingConfigInstance.RealTileSize+g.width > x &&
		g.y*SettingConfigInstance.RealTileSize < y+height &&
		g.y*SettingConfigInstance.RealTileSize+g.height > y {
		return true
	}

	return false
}

func (g *GameObject) Draw(screen *ebiten.Image) {
	g.GameSprite.Draw(screen)

	//DrawRect(screen, g.screenX, g.screenY, g.screenWidth, g.screenHeight, color.RGBA{B: 255, A: 255})

	for _, v := range g.movePosList {
		DrawRect(screen, v.x*SettingConfigInstance.RenderTileSize-CameraInstance.x, v.y*SettingConfigInstance.RenderTileSize-CameraInstance.y, SettingConfigInstance.RenderTileSize, SettingConfigInstance.RenderTileSize, color.RGBA{0, 255, 0, 255})
	}
}
