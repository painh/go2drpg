package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/painh/go2drpg/assetmanager"
	"math"
)

type gameObjectManager struct {
	Width, Height float64
	tiles         []*GameSprite
	objects       []*GameObject
}

func (g *gameObjectManager) Clear() {
	g.tiles = []*GameSprite{}
	g.objects = []*GameObject{}
}

func (g *gameObjectManager) Draw(screen *ebiten.Image) {
	for _, e := range g.tiles {
		e.Draw(screen)
	}

	for _, e := range g.objects {
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

func (g *gameObjectManager) Update(x, y int) {
	objFound := false
	inbound := false

	if x >= ConfigInstance.MapX && y >= ConfigInstance.MapY &&
		x < ConfigInstance.MapWidth && y < ConfigInstance.MapHeight {
		inbound = true
	}

	if InputInstance.LBtnPressed() && inbound {
		if !objFound {
			dx := float64(InputInstance.prevX - x)
			dy := float64(InputInstance.prevY - y)
			CameraInstance.SetXY(CameraInstance.x+dx, CameraInstance.y+dy)
			g.Refresh()
		}
	}

	if inbound && (InputInstance.LBtnClicked() || InputInstance.RBtnClicked()) {
		worldX := float64(int(x+int(CameraInstance.x)) / int(TILE_SIZE))
		worldY := float64(int(y+int(CameraInstance.y)) / int(TILE_SIZE))

		if worldX >= 0 && worldY >= 0 && worldX < GameInstance.mapWidth && worldY < GameInstance.mapHeight {
			for _, e := range g.objects {
				if e.selected {
					e.FindTo(worldX, worldY)
				}
			}

			//GameInstance.cameraToCenter()
			//scripts.StartEvent("slime")
		} else {
			GameInstance.Log.AddWithPrompt("그곳으로 이동 할 수 없습니다.")
		}
	}

	_, wy := ebiten.Wheel()
	if wy != 0 {
		TILE_SIZE = math.Max(1, TILE_SIZE+wy)
		SCALE = TILE_SIZE / SPRITE_PATTERN

		GameInstance.cameraToCenter()

		//CameraInstance.Refresh()
		//g.Refresh()
	}

	for _, e := range g.objects {
		e.Update()
	}
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

func (g *gameObjectManager) GameSpriteAdd(x, y, width, height float64, name string) {
	img := assetmanager.Get(name)
	obj := &GameSprite{x: x, y: y, screenX: 0, screenY: 0, width: width, height: height, selected: false, img: img, op: &ebiten.DrawImageOptions{}}
	obj.Init()
	obj.SetXY(x, y)
	g.tiles = append(g.tiles, obj)
}

func (g *gameObjectManager) GameObjectAdd(x, y, width, height float64, sprName, objName string) {
	img := assetmanager.Get(sprName)
	obj := &GameObject{GameSprite: GameSprite{x: x, y: y, screenX: 0, screenY: 0, width: width, height: height, selected: false, img: img, op: &ebiten.DrawImageOptions{}},
		cdmanager:   CooldownManager{dict: make(map[string]*Cooldown)},
		objName:     objName,
		movePosList: []*TilePos{}}
	//obj := &GameSprite{x, y, screenWidth, screenHeight, false, img, ebiten.DrawImageOptions{}}
	obj.SetXY(x, y)
	obj.SetSize(width, height)

	if objName == ConfigInstance.PlayerObjectName {
		obj.selected = true
	}

	obj.Init()
	g.objects = append(g.objects, obj)
}

func (g *gameObjectManager) CheckGameObjectPosition(x, y, width, height float64, self *GameObject) *GameObject {
	for _, e := range g.objects {
		if e == self {
			continue
		}
		if e.CheckCollision(x, y, width, height) {
			return e
		}
	}
	return nil
}
