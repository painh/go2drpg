package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/painh/go2drpg/assetmanager"
	"sort"
)

type gameObjectManager struct {
	Width, Height float64
	tiles         []*GameSprite
	objects       []*GameObject
	overObjects   []*GameObject
	activeObject  *GameObject
}

func (g *gameObjectManager) Clear() {
	g.tiles = []*GameSprite{}
	g.objects = []*GameObject{}
	g.overObjects = []*GameObject{}
}

func (g *gameObjectManager) Draw(screen *ebiten.Image) {
	for _, e := range g.tiles {
		e.Draw(screen)
	}

	for _, e := range g.objects {
		if g.activeObject == e {
			e.DrawSelected(screen, e.isChar)
		}
		e.Draw(screen)
	}

	for _, e := range g.overObjects {
		e.Draw(screen)
	}
}

func (g *gameObjectManager) Refresh() {
	for _, e := range g.tiles {
		e.Refresh()
	}

	for _, e := range g.objects {
		e.Refresh()
	}

	for _, e := range g.overObjects {
		e.Refresh()
	}
}

func (g *gameObjectManager) selectProcess(x, y float64) bool {
	objFound := false

	if InputInstance.LBtnPressed() {
		for _, e := range g.objects {
			e.selected = false
			if e.clickCheck(float64(x), float64(y)) {
				e.selected = true
				//fmt.Println(e)
				objFound = true
			}
		}
	}

	return objFound
}

func (g *gameObjectManager) Inbound(x, y int) bool {
	if x >= SettingConfigInstance.MapX && y >= SettingConfigInstance.MapY &&
		x < SettingConfigInstance.MapWidth && y < SettingConfigInstance.MapHeight {
		return true
	}

	return false
}

func (g *gameObjectManager) Update(x, y int) {
	objFound := false
	inbound := g.Inbound(x, y)

	if InputInstance.LBtnPressed() && inbound {
		if !objFound {
			dx := float64(InputInstance.prevX - x)
			dy := float64(InputInstance.prevY - y)
			CameraInstance.SetXY(CameraInstance.x+dx, CameraInstance.y+dy)
			g.Refresh()
		}
	}

	if inbound && InputInstance.RBtnClicked() {
		worldX := float64(int(x+int(CameraInstance.x)) / int(SettingConfigInstance.RenderTileSize))
		worldY := float64(int(y+int(CameraInstance.y)) / int(SettingConfigInstance.RenderTileSize))
		g.activeObject = nil

		if worldX >= 0 && worldY >= 0 && worldX < GameInstance.mapWidth && worldY < GameInstance.mapHeight {
			for _, e := range g.objects {
				if e.selected {
					e.FindTo(worldX, worldY)
				}
			}

			//GameInstance.cameraToCenter()
			//scripts.StartEvent("slime")
		} else {
			GameInstance.log.AddWithPrompt("그곳으로 이동 할 수 없습니다.")
		}
	}

	//_, wy := ebiten.Wheel()
	//if wy != 0 {
	//	SettingConfigInstance.RenderTileSize = math.Max(1, SettingConfigInstance.RenderTileSize+wy)
	//	GameInstance.scale = SettingConfigInstance.RenderTileSize / SettingConfigInstance.RealTileSize
	//
	//	GameInstance.cameraToCenter()
	//
	//	//CameraInstance.Refresh()
	//	//g.Refresh()
	//}

	moved := false

	for _, e := range g.objects {
		prevX := e.x
		prevY := e.y
		e.Update()

		if prevX != e.x || prevY != e.y {
			moved = true
		}
	}

	if moved {
		g.SortYPos()
	}

	for _, e := range g.overObjects {
		e.Update()
	}
}

func (g *gameObjectManager) SortYPos() {
	sort.Slice(g.objects, func(i, j int) bool {
		if g.objects[i].y < g.objects[j].y {
			return true
		}
		return false
	})
}

func (g *gameObjectManager) GetSelectedList() []*GameObject {
	var ret = []*GameObject{}

	for _, e := range g.objects {
		if e.selected {
			ret = append(ret, e)
		}
	}

	return ret
}

func (g *gameObjectManager) GameSpriteAdd(x, y, width, height float64, img *assetmanager.ImageResource) {
	obj := &GameSprite{x: x, y: y, screenX: 0, screenY: 0, width: width, height: height, selected: false, img: img, op: &ebiten.DrawImageOptions{}}
	obj.Init()
	obj.SetXY(x, y)
	g.tiles = append(g.tiles, obj)
}

func (g *gameObjectManager) GameObjectAdd(x, y, width, height float64, img *assetmanager.ImageResource, objName string, zindex float64, isOverObject bool, isChar bool) *GameObject {
	obj := &GameObject{GameSprite: GameSprite{x: x, y: y, screenX: 0, screenY: 0, width: width, height: height, selected: false, img: img, op: &ebiten.DrawImageOptions{}},
		cdmanager:   CooldownManager{dict: make(map[string]*Cooldown)},
		objName:     objName,
		movePosList: []*TilePos{},
		isChar:      isChar}
	//obj := &GameSprite{x, y, screenWidth, screenHeight, false, img, ebiten.DrawImageOptions{}}
	obj.SetXY(x, y)
	obj.SetSize(width, height)
	obj.zindex = zindex

	if objName == SettingConfigInstance.PlayerObjectName {
		obj.selected = true
	}

	obj.Init()

	if isOverObject {
		g.overObjects = append(g.overObjects, obj)
	} else {
		g.objects = append(g.objects, obj)
	}

	return obj
}

func (g *gameObjectManager) CheckGameObjectPosition(x, y, z, width, height float64, self *GameObject) *GameObject {
	for _, e := range g.objects {
		if e == self {
			continue
		}

		if z != e.zindex {
			continue
		}
		if e.CheckCollision(x, y, width, height) {
			return e
		}
	}
	return nil
}

func (g *gameObjectManager) SetActiveObject(obj *GameObject) {
	g.activeObject = obj

}
