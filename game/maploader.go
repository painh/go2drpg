package game

import (
	"github.com/fardog/tmx"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/painh/go2drpg/assetmanager"
	"log"
	"math"
)

type MapLoader struct {
	Width    int
	Height   int
	tmxMap   *tmx.Map
	tilesets []tmx.TileSet
}

func (m *MapLoader) GetTileImage(n tmx.GlobalID) *assetmanager.ImageResource {
	for k, v := range m.tilesets {
		max := tmx.GlobalID(math.MaxInt32)
		if k != len(m.tilesets)-1 {
			max = m.tilesets[k+1].FirstGlobalID
		}

		if n >= v.FirstGlobalID && n < max {
			tileset := m.tilesets[k]
			tilenum := n - (tileset.FirstGlobalID)
			img := assetmanager.GetWithName(tileset.Name, int(tilenum))
			if n != 0 && img == nil {
				log.Fatalf("img not found %v, %v, %v", tileset.Name, n, tilenum)
			}

			return img
		}
	}

	return nil
}

func (m *MapLoader) Load(filename string, g *Game) {
	file, err := ebitenutil.OpenFile(SettingConfigInstance.WorkFolder + filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer file.Close()

	m.tmxMap, err = tmx.Decode(file)
	if err != nil {
		log.Fatal(err.Error())
	}

	if m.tmxMap == nil {
		log.Fatal("map was nil")
	}

	m.tilesets = []tmx.TileSet{}

	for _, v := range m.tmxMap.TileSets {
		//v.Image
		file, err := ebitenutil.OpenFile(SettingConfigInstance.WorkFolder + v.Source)
		if err != nil {
			log.Fatal(err.Error())
		}

		tileset, err := tmx.DecodeTileset(file)
		if err != nil {
			log.Fatal(err.Error())
		}

		file.Close()

		v.Name = tileset.Name
		m.tilesets = append(m.tilesets, v)

		assetmanager.Load(SettingConfigInstance.WorkFolder+tileset.Image.Source, tileset.Name)
		assetmanager.MakePatternImages(tileset.Name, tileset.TileWidth, tileset.TileHeight)
	}
	//m.tmxMap.TileSets[0]

	g.gameObjectManager.Clear()

	for _, v := range m.tmxMap.Layers {

		trs, err := v.TileGlobalRefs()
		if err != nil {
			log.Fatal(err)
		}

		g.mapWidth = float64(m.tmxMap.Width)
		g.mapHeight = float64(m.tmxMap.Height)

		if l, e := len(trs), m.tmxMap.Width*m.tmxMap.Height; l != e {
			log.Fatalf("expected tiles of length %v, got %v", e, l)
		}

		tds, err := v.TileDefs(m.tmxMap.TileSets)
		if err != nil {
			log.Fatal(err)
		}

		cnt := 0

		for y := 0; y < v.Height; y++ {
			for x := 0; x < v.Width; x++ {
				//"base_"+strconv.Itoa())
				image := m.GetTileImage(tmx.GlobalID(tds[cnt].GlobalID))
				g.gameObjectManager.GameSpriteAdd(float64(x), float64(y), SettingConfigInstance.RealTileSize, SettingConfigInstance.RealTileSize, image)
				cnt++
			}
		}
	}

	if objects := m.tmxMap.ObjectGroupWithName("obj"); objects != nil {

		for _, v := range objects.Objects {
			x := v.X / SettingConfigInstance.RealTileSize
			y := v.Y / SettingConfigInstance.RealTileSize

			if v.GlobalID != 0 { //GlobalID가 없다면, 이미지가 선택되지 않은 충돌 사각형으로 생각하여 원본을 씀. Tiled가 오브젝트일때는 y좌표를 + 1해서 주는 방식이라 후처리를 해야함.
				y--
			}

			image := m.GetTileImage(tmx.GlobalID(v.GlobalID))

			g.gameObjectManager.GameObjectAdd(float64(x), float64(y), v.Width, v.Height, image, v.Name, 0, false, false)
		}
	} else {
		log.Fatal("cant find object layer : name {obj}")
	}

	// char은 48x96이므로 특수 처리를 함
	if objects := m.tmxMap.ObjectGroupWithName("char"); objects != nil {

		for _, v := range objects.Objects {
			x := v.X / SettingConfigInstance.RealTileSize
			y := v.Y / SettingConfigInstance.RealTileSize

			if v.GlobalID != 0 { //GlobalID가 없다면, 이미지가 선택되지 않은 충돌 사각형으로 생각하여 원본을 씀. Tiled가 오브젝트일때는 y좌표를 + 1해서 주는 방식이라 후처리를 해야함.
				y--
			}

			image := m.GetTileImage(tmx.GlobalID(v.GlobalID))

			obj := g.gameObjectManager.GameObjectAdd(float64(x), float64(y), v.Width, v.Height/2, image, v.Name, 0, false, true)
			obj.offsetY = -1
			obj.Refresh()
		}
	} else {
		log.Fatal("cant find object layer : name {char}")
	}

	if objects := m.tmxMap.ObjectGroupWithName("overchar"); objects != nil {

		for _, v := range objects.Objects {
			x := v.X / SettingConfigInstance.RealTileSize
			y := v.Y / SettingConfigInstance.RealTileSize

			if v.GlobalID != 0 { //GlobalID가 없다면, 이미지가 선택되지 않은 충돌 사각형으로 생각하여 원본을 씀. Tiled가 오브젝트일때는 y좌표를 + 1해서 주는 방식이라 후처리를 해야함.
				y--
			}

			image := m.GetTileImage(tmx.GlobalID(v.GlobalID))

			g.gameObjectManager.GameObjectAdd(float64(x), float64(y), v.Width, v.Height, image, v.Name, 1, true, false)
		}
	} else {
		log.Fatal("cant find object layer : name {overchar}")
	}

	m.Width = m.tmxMap.Width
	m.Height = m.tmxMap.Height
}
