package assetmanager

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
	"strconv"
)

type ImageResource struct {
	Img      *ebiten.Image
	Filename string
	Pattern  bool
}

type imageManager struct {
	dict map[string]*ImageResource
}

var Instance imageManager

func Load(filename string, name string) *ImageResource {
	if v, found := Instance.dict[filename]; found {
		return v
	}

	res := &ImageResource{Pattern: false, Filename: filename}

	var err error
	res.Img, _, err = ebitenutil.NewImageFromFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	Instance.dict[name] = res

	return res
}

func makePatternImage(origin *ImageResource, patternX, patternY, patternWidth, patternHeight int) *ImageResource {
	res := &ImageResource{Pattern: true, Filename: origin.Filename}
	img := origin.Img.SubImage(image.Rect(patternX, patternY, patternX+patternWidth, patternY+patternHeight)).(*ebiten.Image)
	res.Img = img

	return res
}

func Get(name string) *ImageResource {
	v, found := Instance.dict[name]
	if !found {
		log.Fatal(name + " was not exist")
	}

	return v
}

func MakePatternImages(name string, patternWidth, patternHeight int) {
	origin := Get(name)

	width, height := origin.Img.Size()
	cnt := 0

	for y := 0; y < height; y += patternHeight {
		for x := 0; x < width; x += patternWidth {
			img := makePatternImage(origin, x, y, patternWidth, patternHeight)
			Instance.dict[name+"_"+strconv.Itoa(cnt)] = img
			cnt++
		}
	}
}

func init() {
	Instance.dict = make(map[string]*ImageResource)
}
