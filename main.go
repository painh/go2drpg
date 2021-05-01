// Copyright 2017 The Ebiten Authors
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

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/painh/go2drpg/game"
	_ "image/png"
	"log"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

func main() {
	game.ConfigInstance.Load("config.json")

	g := game.NewGame(game.ConfigInstance.Window_width, game.ConfigInstance.Window_height)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}