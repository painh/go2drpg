package game

type Camera struct {
	x, y             float64
	originX, originY float64
}

var CameraInstance = Camera{}

func (c *Camera) SetXY(x, y float64) {
	c.x = x
	c.y = y
	c.originX = x * SCALE
	c.originY = y * SCALE
}

func (c *Camera) Refresh() {
	c.x = c.originX / SCALE
	c.y = c.originY / SCALE
}
