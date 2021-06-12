package game

type Camera struct {
	x, y             float64
	originX, originY float64
}

var CameraInstance = Camera{}

func (c *Camera) SetXY(x, y float64) {
	c.x = x
	c.y = y
	c.originX = x * GameInstance.scale
	c.originY = y * GameInstance.scale
}

func (c *Camera) Refresh() {
	c.x = c.originX / GameInstance.scale
	c.y = c.originY / GameInstance.scale
}
