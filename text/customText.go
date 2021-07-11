// from ebiten's text
package text

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"math"
	"sync"
)

type RichTextElement struct {
	tag     string
	value   string
	inner   string
	opentag bool
}

var errInvalidFormat = errors.New("invalid format")

// https://stackoverflow.com/questions/54197913/parse-hex-string-to-image-color
func ParseHexColorFast(s string) (c color.RGBA, err error) {
	c.A = 0xff

	if s[0] != '#' {
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errInvalidFormat
	}
	return
}

func SimpleRichParser(fulltext string) []*RichTextElement {
	list := []*RichTextElement{}
	root := &RichTextElement{"root", "0", "", true}
	curTagElement := root
	list = append(list, root) // Push
	text := []rune(fulltext)
	for i := 0; i < len(text); i++ {
		c := text[i]
		if string(c) == "<" && i < len(text)-1 {
			tag := ""
			value := ""
			curText := ""
			if string(text[i+1]) == "/" {
				for i = i + 1; i < len(text); i++ {
					c2 := string(text[i])
					if c2 == ">" {
						prevTag := list[len(list)-2]
						curTagElement = &RichTextElement{tag: prevTag.tag, value: prevTag.value, opentag: false}
						list = append(list, curTagElement) // Push
						break
					}
				}
			} else {
				for i = i + 1; i < len(text); i++ {
					c2 := string(text[i])
					if c2 == "=" {
						tag = curText
						curText = ""
						continue
					}

					if c2 == ">" {
						value = curText
						curTagElement = &RichTextElement{tag: tag, value: value, opentag: true}
						list = append(list, curTagElement) // Push
						break
					}

					curText += c2
				}
			}
		} else {
			curTagElement.inner += string(c)
		}
	}

	return list
}

var glyphAdvanceCache = map[font.Face]map[rune]fixed.Int26_6{}

func glyphAdvance(face font.Face, r rune) fixed.Int26_6 {
	m, ok := glyphAdvanceCache[face]
	if !ok {
		m = map[rune]fixed.Int26_6{}
		glyphAdvanceCache[face] = m
	}

	a, ok := m[r]
	if !ok {
		a, _ = face.GlyphAdvance(r)
		m[r] = a
	}

	return a
}

var (
	monotonicClock int64
)

func now() int64 {
	return monotonicClock
}

func init() {
	//hooks.AppendHookOnBeforeUpdate(func() error {
	//	monotonicClock++
	//	return nil
	//})
}

func Update() {
	monotonicClock++
}

func fixed26_6ToFloat64(x fixed.Int26_6) float64 {
	return float64(x>>6) + float64(x&((1<<6)-1))/float64(1<<6)
}

func drawGlyph(dst *ebiten.Image, face font.Face, r rune, img *ebiten.Image, x, y fixed.Int26_6, clr ebiten.ColorM) {
	if img == nil {
		return
	}

	b := getGlyphBounds(face, r)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64((x+b.Min.X)>>6), float64((y+b.Min.Y)>>6))
	op.ColorM = clr
	dst.DrawImage(img, op)
}

var (
	glyphBoundsCache = map[font.Face]map[rune]fixed.Rectangle26_6{}
)

func getGlyphBounds(face font.Face, r rune) fixed.Rectangle26_6 {
	if _, ok := glyphBoundsCache[face]; !ok {
		glyphBoundsCache[face] = map[rune]fixed.Rectangle26_6{}
	}
	if b, ok := glyphBoundsCache[face][r]; ok {
		return b
	}
	b, _, _ := face.GlyphBounds(r)
	glyphBoundsCache[face][r] = b
	return b
}

type glyphImageCacheEntry struct {
	image *ebiten.Image
	atime int64
}

var (
	glyphImageCache = map[font.Face]map[rune]*glyphImageCacheEntry{}
)

func getGlyphImage(face font.Face, r rune) *ebiten.Image {
	if _, ok := glyphImageCache[face]; !ok {
		glyphImageCache[face] = map[rune]*glyphImageCacheEntry{}
	}

	if e, ok := glyphImageCache[face][r]; ok {
		e.atime = now()
		return e.image
	}

	b := getGlyphBounds(face, r)
	w, h := (b.Max.X - b.Min.X).Ceil(), (b.Max.Y - b.Min.Y).Ceil()
	if w == 0 || h == 0 {
		glyphImageCache[face][r] = &glyphImageCacheEntry{
			image: nil,
			atime: now(),
		}
		return nil
	}

	if b.Min.X&((1<<6)-1) != 0 {
		w++
	}
	if b.Min.Y&((1<<6)-1) != 0 {
		h++
	}
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))

	d := font.Drawer{
		Dst:  rgba,
		Src:  image.White,
		Face: face,
	}
	x, y := -b.Min.X, -b.Min.Y
	x, y = fixed.I(x.Ceil()), fixed.I(y.Ceil())
	d.Dot = fixed.Point26_6{X: x, Y: y}
	d.DrawString(string(r))

	img := ebiten.NewImageFromImage(rgba)
	if _, ok := glyphImageCache[face][r]; !ok {
		glyphImageCache[face][r] = &glyphImageCacheEntry{
			image: img,
			atime: now(),
		}
	}

	return img
}

var textM sync.Mutex

// Draw draws a given text on a given destination image dst.
//
// face is the font for text rendering.
// (x, y) represents a 'dot' (period) position.
// This means that if the given text consisted of a single character ".",
// it would be positioned at the given position (x, y).
// Be careful that this doesn't represent left-upper corner position.
//
// clr is the color for text rendering.
//
// If you want to adjust the position of the text, these functions are useful:
//
//     * text.BoundString:                     the rendered bounds of the given text.
//     * golang.org/x/image/font.Face.Metrics: the metrics of the face.
//
// The '\n' newline character puts the following text on the next line.
// Line height is based on Metrics().Height of the font.
//
// Glyphs used for rendering are cached in least-recently-used way.
// Then old glyphs might be evicted from the cache.
// As the cache capacity has limit, it is not guaranteed that all the glyphs for runes given at Draw are cached.
// The cache is shared with CacheGlyphs.
//
// It is OK to call Draw with a same text and a same face at every frame in terms of performance.
//
// Draw and CacheGlyphs are implemented like this:
//
//     Draw        = Create glyphs by `(*ebiten.Image).ReplacePixels` and put them into the cache if necessary
//                 + Draw them onto the destination by `(*ebiten.Image).DrawImage`
//     CacheGlyphs = Create glyphs by `(*ebiten.Image).ReplacePixels` and put them into the cache if necessary
//
// Be careful that the passed font face is held by this package and is never released.
// This is a known issue (#498).
//
// Draw is concurrent-safe.
func Draw(dst *ebiten.Image, fullText string, face font.Face, x, y int, clr color.Color, justBound bool, wordWrap float64) image.Rectangle {
	textM.Lock()
	defer textM.Unlock()

	cr, cg, cb, ca := clr.RGBA()
	//if ca == 0 {
	//	return
	//}

	var colorm ebiten.ColorM
	//if justBound == false {
	//	colorm.Scale(float64(cr)/float64(ca), float64(cg)/float64(ca), float64(cb)/float64(ca), float64(ca)/0xffff)
	//}

	fx, fy, ww := fixed.I(x), fixed.I(y), fixed.I(int(wordWrap))
	if ww > 0 {
		ww += fx
	}
	prevR := rune(-1)

	list := SimpleRichParser(fullText)
	faceHeight := face.Metrics().Height
	fy += faceHeight
	var bounds fixed.Rectangle26_6

	for _, text := range list {
		if justBound == false {
			if text.tag == "color" {
				c, _ := ParseHexColorFast(text.value)
				cr, cg, cb, ca = c.RGBA()
			} else {
				cr, cg, cb, ca = clr.RGBA()
			}
			colorm.Reset()
			colorm.Scale(float64(cr)/float64(ca), float64(cg)/float64(ca), float64(cb)/float64(ca), float64(ca)/0xffff)
			//fmt.Printf("color : %v %v %v %v text:%v\n", cr, cg, cb, ca, text.inner)
		}
		for _, r := range text.inner {
			if prevR >= 0 {
				fx += face.Kern(prevR, r)
			}

			b := getGlyphBounds(face, r)
			b.Min.X += fx
			b.Max.X += fx
			b.Min.Y += fy
			b.Max.Y += fy
			bounds = bounds.Union(b)

			if r == '\n' {
				fx = fixed.I(x)
				fy += faceHeight
				prevR = rune(-1)
				continue
			}

			dx := b.Max.X - b.Min.X
			if ww > 0 && fx+dx >= ww {
				fx = fixed.I(x)
				fy += faceHeight
			}

			if justBound == false {
				img := getGlyphImage(face, r)
				drawGlyph(dst, face, r, img, fx, fy, colorm)
			}
			fx += glyphAdvance(face, r)

			prevR = r
		}
	}

	// cacheSoftLimit indicates the soft limit of the number of glyphs in the cache.
	// If the number of glyphs exceeds this soft limits, old glyphs are removed.
	// Even after clearning up the cache, the number of glyphs might still exceeds the soft limit, but
	// this is fine.
	const cacheSoftLimit = 512

	// Clean up the cache.
	if len(glyphImageCache[face]) > cacheSoftLimit {
		for r, e := range glyphImageCache[face] {
			// 60 is an arbitrary number.
			if e.atime < now()-60 {
				delete(glyphImageCache[face], r)
			}
		}
	}

	return image.Rect(
		int(math.Floor(fixed26_6ToFloat64(bounds.Min.X))),
		int(math.Floor(fixed26_6ToFloat64(bounds.Min.Y))),
		int(math.Ceil(fixed26_6ToFloat64(bounds.Max.X))),
		int(math.Ceil(fixed26_6ToFloat64(bounds.Max.Y))),
	)
}

func BoundString(face font.Face, text string) image.Rectangle {
	return Draw(nil, text, face, 0, 0, color.RGBA{0, 0, 0, 0}, true, 0)
}

// CacheGlyphs precaches the glyphs for the given text and the given font face into the cache.
//
// Glyphs used for rendering are cached in least-recently-used way.
// Then old glyphs might be evicted from the cache.
// As the cache capacity has limit, it is not guaranteed that all the glyphs for runes given at CacheGlyphs are cached.
// The cache is shared with Draw.
//
// Draw and CacheGlyphs are implemented like this:
//
//     Draw        = Create glyphs by `(*ebiten.Image).ReplacePixels` and put them into the cache if necessary
//                 + Draw them onto the destination by `(*ebiten.Image).DrawImage`
//     CacheGlyphs = Create glyphs by `(*ebiten.Image).ReplacePixels` and put them into the cache if necessary
//
// Draw automatically creates and caches necessary glyphs, so usually you don't have to call CacheGlyphs
// explicitly. However, for example, when you call Draw for each rune of one big text, Draw tries to create the glyph
// cache and render it for each rune. This is very inefficient because creating a glyph image and rendering it are
// different operations (`(*ebiten.Image).ReplacePixels` and `(*ebiten.Image).DrawImage`) and can never be merged as
// one draw call. CacheGlyphs creates necessary glyphs without rendering them so that these operations are likely
// merged into one draw call regardless of the size of the text.
//
// If a rune's glyph is already cached, CacheGlyphs does nothing for the rune.
func CacheGlyphs(face font.Face, text string) {
	textM.Lock()
	defer textM.Unlock()

	for _, r := range text {
		getGlyphImage(face, r)
	}
}
