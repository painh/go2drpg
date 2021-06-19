package game

type KeywordManager struct {
	dict map[string]bool
}

func (i *KeywordManager) Init() {
	i.dict = map[string]bool{}
}

func (i *KeywordManager) Add(key string) {
	if i.CheckKeyword(key) {
		return
	}
	GameInstance.log.AddWithPrompt("키워드 ", key, "(이)가 추가 되었습니다.")
	i.dict[key] = true
}

func (i *KeywordManager) CheckKeyword(key string) bool {
	_, ok := i.dict[key]
	return ok
}
