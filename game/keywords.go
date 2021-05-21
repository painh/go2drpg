package game

import "github.com/hajimehoshi/ebiten/v2"

type KeyWords struct {
	id   int
	name string
	img  ebiten.Image
	desc string
}

type KeyWordsManager struct {
	list []KeyWords
}
