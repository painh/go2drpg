package game

import (
	"github.com/beefsack/go-astar"
)

type TilePos struct {
	x float64
	y float64
}

type TileManager struct {
	dict           map[float64]map[float64]*TilePos
	x1, y1, x2, y2 float64
	currentObject  *GameObject
}

func (t *TileManager) FindTo(x1, y1, x2, y2 float64, currentObject *GameObject) ([]astar.Pather, float64, bool) {
	t.x1 = x1
	t.y1 = y1
	t.x2 = x2
	t.y2 = y2
	t.currentObject = currentObject

	from := tileManagerInstance.Get(x1, y1)
	to := tileManagerInstance.Get(x2, y2)
	path, distance, found := astar.Path(from, to)
	//fmt.Println(path)
	//fmt.Println(distance)
	//fmt.Println(found)

	return path, distance, found
}

func (t *TileManager) Init() {
	t.dict = make(map[float64]map[float64]*TilePos)
}

func (t *TileManager) Get(x, y float64) *TilePos {
	_, found := t.dict[x]

	if !found {
		t.dict[x] = make(map[float64]*TilePos)
	}

	v, found := t.dict[x][y]

	if found {
		return v
	}

	t.dict[x][y] = &TilePos{x, y}
	return t.dict[x][y]
}

func (t *TileManager) IsToPos(x, y float64) bool {
	if t.x2 == x && t.y2 == y {
		return true
	}

	return false
}

var tileManagerInstance TileManager = TileManager{}

func init() {
	tileManagerInstance.Init()
}

func (t *TilePos) ManhattanDistance(to *TilePos) float64 {
	return ManhattanDistance(t.x, t.y, to.x, to.y)
}

func (t *TilePos) PathNeighbors() []astar.Pather {

	var ret []astar.Pather
	const step = 1

	list := [][]float64{
		{t.x - step, t.y},
		{t.x - step, t.y - step},
		{t.x, t.y - step},
		{t.x + step, t.y - step},
		{t.x + step, t.y},
		{t.x + step, t.y + step},
		{t.x, t.y + step},
		{t.x - step, t.y + step},
	}

	for _, v := range list {
		distnace := ManhattanDistance(tileManagerInstance.x1, tileManagerInstance.y1, v[0], v[1])
		if distnace > 100 {
			//log.Println("too far", tileManagerInstance.x1, tileManagerInstance.y1, v[0], v[1], distnace)
			continue
		}

		if tileManagerInstance.IsToPos(v[0], v[1]) { //목적지는 언제나 갈수 있음.
			tile := tileManagerInstance.Get(v[0], v[1])
			ret = append(ret, tile)
			continue
		}

		obj := GameInstance.gameObjectManager.CheckGameObjectPosition(v[0]*SPRITE_PATTERN, v[1]*SPRITE_PATTERN, tileManagerInstance.currentObject.width, tileManagerInstance.currentObject.height, tileManagerInstance.currentObject)
		if obj == nil {
			tile := tileManagerInstance.Get(v[0], v[1])
			ret = append(ret, tile)
		}
	}

	//fmt.Println(ret)

	return ret
}

func (t *TilePos) PathNeighborCost(to astar.Pather) float64 {
	return TILE_SIZE //이동비용은 항상 1
}

func (t *TilePos) PathEstimatedCost(to astar.Pather) float64 {
	return t.ManhattanDistance(to.(*TilePos))
}
