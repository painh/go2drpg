package game

import (
	"fmt"
	"math"
)

type DoubleClick struct {
	LastX         int
	LastY         int
	LastClickTS   int64
	DoubleClicked bool
}

var DoubleClickInstance = DoubleClick{0, 0, 0, false}

func (d *DoubleClick) Update(x, y int, pressed bool) {
	d.DoubleClicked = false
	if !pressed {
		return
	}

	now := makeTimestamp()
	fmt.Println(now-d.LastClickTS, ConfigInstance.Doubclick_ts_margin)

	if math.Abs(float64(x-d.LastX)) < float64(ConfigInstance.Doubclick_pixel_margin) &&
		math.Abs(float64(y-d.LastY)) < float64(ConfigInstance.Doubclick_pixel_margin) &&
		now-d.LastClickTS < ConfigInstance.Doubclick_ts_margin {
		d.DoubleClicked = true
	}

	d.LastX = x
	d.LastY = y
	d.LastClickTS = now
}

func (d *DoubleClick) DBClick() bool {
	return d.DoubleClicked
}
