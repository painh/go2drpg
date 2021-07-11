package game

import (
	"fmt"
	"image/color"
)

type KeywordManager struct {
	dict          map[string]bool
	activeKeyword string
}

func (i *KeywordManager) Init() {
	i.dict = map[string]bool{}
}

func (i *KeywordManager) Add(key string) {
	if i.CheckKeyword(key) {
		return
	}
	GameInstance.log.AddWithPrompt(color.RGBA{0, 255, 0, 255},"키워드 ", key, "(이)가 추가 되었습니다.")
	i.dict[key] = true
}

func (i *KeywordManager) CheckKeyword(key string) bool {
	_, ok := i.dict[key]
	return ok
}

func (i *KeywordManager) ActiveKeyword(keyword string) {
	if i.activeKeyword == keyword {
		fmt.Println("keyword already activated : ", keyword)
	}
	i.activeKeyword = keyword
}
