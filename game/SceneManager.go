package game

type SceneManager struct {
	scenename              string
	scene                  []ScriptActionInterface
	cursor                 int
	isOver                 bool
	person                 bool
	nonexclusive           bool
	condition              []interface{}
	invalidKeywordResponse string
}

func (s *SceneManager) Start() {
	s.isOver = false
	s.cursor = 0
}

func (s *SceneManager) Update() bool {
	if s.cursor >= len(s.scene) {
		s.isOver = true
	}

	if s.isOver {
		return false
	}

	currentScript := s.scene[s.cursor]
	currentScript.Run()

	s.cursor++

	return true

}

func (s *SceneManager) CheckCondition(keyword string) bool {
	for _, v := range s.condition {
		for k2, v2 := range v.(map[interface{}]interface{}) {
			switch k2 {
			case "keyword":
				if keyword != v2 {
					return false
				}
			case "switch":
				if v, ok := GameInstance.gameSwitchManager.CheckSwitch(v2.(string)); !ok || v == false {
					return false
				}

			case "location":
				if GameInstance.player.currentLocationName != k2.(string) {
					return false
				}
			}

		}
	}

	return true
}

func (s *SceneManager) GetConditionKeyCount(key string) int {
	ret := 0

	for _, v := range s.condition {
		for k2, _ := range v.(map[interface{}]interface{}) {
			if k2 == key {
				ret++
			}
		}
	}

	return ret
}

func (s *SceneManager) FindAvailableKeywordScene(key string) bool {
	keywordCnt := s.GetConditionKeyCount("keyword")

	if keywordCnt == 0 {
		return false
	}

	if s.CheckCondition(key) == false {
		return false
	}

	return true
}
