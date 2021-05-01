package game

import (
	"github.com/beefsack/go-astar"
)
import "math"

type TilePos struct {
	x float64
	y float64
}

type TileManager struct {
	dict map[float64]map[float64]*TilePos
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

var tileManager TileManager = TileManager{}

func init() {
	tileManager.Init()
}

var TILE_SIZE float64 = 32

const SPRITE_PATTERN = float64(16)

var SCALE float64 = TILE_SIZE / SPRITE_PATTERN

func (t *TilePos) ManhattanDistance(to *TilePos) float64 {
	return math.Abs(t.x-to.x) + math.Abs(t.y-to.y)
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
		if !GameInstance.gameObjectManager.CheckGameObjectPosition(v[0], v[1]) {
			tile := tileManager.Get(v[0], v[1])
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
