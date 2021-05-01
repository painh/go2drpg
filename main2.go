// Copyright 2016 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build example

package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

const sampleText = `The quick brown fox jumps over the lazy dog.`

var (
	mplusNormalFont font.Face
	mplusBigFont    font.Face
	jaKanjis        = []rune{}
)

func init() {
	// table is the list of Japanese Kanji characters in a part of JIS X 0208.
	const table = `
가가가가가나나나나
`
	for _, c := range table {
		if c == '\n' {
			continue
		}
		jaKanjis = append(jaKanjis, c)
	}
}

func init() {
	dat, err := ioutil.ReadFile("SCDream3.otf")

	tt, err := opentype.Parse(dat)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Game struct {
	counter        int
	kanjiText      []rune
	kanjiTextColor color.RGBA
}

func (g *Game) Update() error {
	// Change the text color for each second.
	if g.counter%ebiten.MaxTPS() == 0 {
		g.kanjiText = nil
		for j := 0; j < 4; j++ {
			for i := 0; i < 8; i++ {
				g.kanjiText = append(g.kanjiText, jaKanjis[rand.Intn(len(jaKanjis))])
			}
			g.kanjiText = append(g.kanjiText, '\n')
		}

		g.kanjiTextColor.R = 0x80 + uint8(rand.Intn(0x7f))
		g.kanjiTextColor.G = 0x80 + uint8(rand.Intn(0x7f))
		g.kanjiTextColor.B = 0x80 + uint8(rand.Intn(0x7f))
		g.kanjiTextColor.A = 0xff
	}
	g.counter++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	const x = 20

	// Draw info
	msg := fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS())
	text.Draw(screen, msg, mplusNormalFont, x, 40, color.White)

	// Draw the sample text
	text.Draw(screen, sampleText, mplusNormalFont, x, 80, color.White)

	// Draw Kanji text lines
	for i, line := range strings.Split(string(g.kanjiText), "\n") {
		text.Draw(screen, line, mplusBigFont, x, 160+54*i, g.kanjiTextColor)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Font (Ebiten Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
