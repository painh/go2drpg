package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/painh/go2drpg/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"log"
)

type FontText struct {
	fontface font.Face
}

func (f *FontText) FontHeight() int {
	b, _, _ := f.fontface.GlyphBounds('M')
	return (b.Max.Y - b.Min.Y).Ceil()
}

func (f *FontText) LoadFont(filename string, size int) {
	dat, err := ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(dat)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	f.fontface, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

func (f *FontText) BoundString(str string) image.Rectangle {
	return text.BoundString(f.fontface, str)
}

func (f *FontText) Draw(dst *ebiten.Image, str string, x, y int, rgba color.RGBA) image.Rectangle {
	return text.Draw(dst, str, f.fontface, int(x), y, rgba, false, 0)
}

func (f *FontText) DrawWithWW(dst *ebiten.Image, str string, x, y int, rgba color.RGBA, jb bool, ww float64) image.Rectangle {
	return text.Draw(dst, str, f.fontface, int(x), y, rgba, jb, ww)
}

func (f *FontText) DrawTextInBox(dst *ebiten.Image, str string, x, y float64) float64 {
	rect := text.BoundString(f.fontface, str)
	ebitenutil.DrawRect(dst, x, y, float64(rect.Dx()), float64(rect.Dy()), color.Black)

	text.Draw(dst, str, f.fontface, int(x), int(y), color.White, false, 0)

	return float64(rect.Dy())
}

var defaultFontInstance = FontText{}
