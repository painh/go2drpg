package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
)

type Input struct {
	x             int
	y             int
	LastX         int
	LastY         int
	LastClickTS   int64
	DoubleClicked bool
	pressed       bool
	clicked       bool
	firstclick    bool
	prevX         int
	prevY         int
	lbtnJustUp    bool
}

var InputInstance = Input{}

func (d *Input) Update() {
	prevPressed := d.pressed
	d.pressed = false
	d.clicked = false
	d.lbtnJustUp = false

	ids := ebiten.TouchIDs()
	if ids == nil {
		d.prevX = d.x
		d.prevY = d.y
		d.x, d.y = ebiten.CursorPosition()
		d.pressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
		d.clicked = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	} else {
		//fmt.Println(ids)
		d.pressed = true

		d.clicked = len(inpututil.JustPressedTouchIDs()) > 0

		d.prevX = d.x
		d.prevY = d.y

		d.x, d.y = ebiten.TouchPosition(ids[0])

		if d.clicked {
			d.prevX = d.x
			d.prevY = d.y
		}
	}

	d.DoubleClicked = false
	if !d.clicked {
		return
	}

	now := makeTimestamp()
	//fmt.Println(now-d.LastClickTS, SettingConfigInstance.DoubleClickTsMargin)

	if math.Abs(float64(d.x-d.LastX)) < float64(SettingConfigInstance.DoubleClickPixelMargin) &&
		math.Abs(float64(d.y-d.LastY)) < float64(SettingConfigInstance.DoubleClickPixelMargin) &&
		now-d.LastClickTS < SettingConfigInstance.DoubleClickTsMargin {
		d.DoubleClicked = true
	}

	d.LastX = d.x
	d.LastY = d.y
	d.LastClickTS = now

	if d.pressed == false && prevPressed {
		d.lbtnJustUp = true
	}
}

func (d *Input) DBClick() bool {
	return d.DoubleClicked
}

func (d *Input) LBtnPressed() bool {
	return d.pressed
}

func (d *Input) LBtnJustUp() bool {
	return d.lbtnJustUp
}

func (d *Input) LBtnClicked() bool {
	return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || d.clicked
}

func (d *Input) RBtnClicked() bool {
	ret := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) || d.DBClick()

	if ret {
		//fmt.Println(inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight), d.DBClick())
	}

	return ret
}
