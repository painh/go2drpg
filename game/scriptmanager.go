package game

import "image/color"

type ScriptManager struct {
	scenes            []*SceneManager
	lastActiveScene   *SceneManager
	activeScriptStack []*SceneManager
	lastObjectName    string
}

func (s *ScriptManager) GetInvalidKeywordResponse() string {
	if s.lastActiveScene == nil {
		return "무슨 말인지 모르겠군요."
	}

	if s.lastActiveScene.invalidKeywordResponse == "" {
		return "무슨 말인지 모르겠군요."
	}

	return s.lastActiveScene.invalidKeywordResponse
}

func (s *ScriptManager) Init() {
	s.scenes = []*SceneManager{}
	s.activeScriptStack = []*SceneManager{}
	s.lastActiveScene = nil
}

func (s *ScriptManager) Update() {
	if s.lastActiveScene == nil {
		return
	}

	if s.lastActiveScene.Update() == false {
		scene := s.PopScene()

		if scene != nil {
			s.lastActiveScene = scene
		}
	}
}

//
//func (s *ScriptManager) Add(scenename string, scr *SceneManager) {
//	s.scenes[scenename] = scr
//	s.scenes[scenename].scenename = scenename
//}

func (s *ScriptManager) FindScript(scenename string) *SceneManager {
	for _, v := range s.scenes {
		if v.scenename == scenename {
			return v
		}
	}

	return nil
}

func (s *ScriptManager) FindAvailableKeywordScene(keyword string) bool {
	for _, v := range s.scenes {
		if v.FindAvailableKeywordScene(keyword) {
			return true
		}
	}

	return false
}

func (s *ScriptManager) RunCurrentObject() bool {
	return s.RunObjectScript(s.lastObjectName)
}

func (s *ScriptManager) RunObjectScript(objName string) bool {
	//TODO : keyword 실행과 합쳐야함
	for _, v := range s.scenes {
		if objName == v.scenename {
			if v.CheckCondition(GameInstance.keywordManager.activeKeyword) {
				s.PushScene()
				s.ActiveScene(v)

				if s.lastActiveScene.person && s.lastObjectName == "" {
					s.lastObjectName = objName
					GameInstance.log.AddWithPrompt(color.RGBA{0, 255, 0, 255},objName, "(와)과 대화를 시작합니다.\n")
				}

				if v.nonexclusive == false {
					return true
				}
			}
		}
	}

	return true
}

func (s *ScriptManager) ActiveScene(scene *SceneManager) bool {
	s.lastActiveScene = scene
	s.lastActiveScene.Start()

	return true
}

func (s *ScriptManager) PushScene() {
	if s.lastActiveScene == nil {
		return
	}
	s.activeScriptStack = append(s.activeScriptStack, s.lastActiveScene)
}

func (s *ScriptManager) PopScene() *SceneManager {
	if len(s.activeScriptStack) < 1 {
		return nil
	}

	lastIdx := len(s.activeScriptStack) - 1
	ret := s.activeScriptStack[lastIdx]
	s.activeScriptStack = s.activeScriptStack[:lastIdx]

	return ret
}

//
//func (s *ScriptManager) RunKeywordScript(keyword string) bool {
//	for _, v := range s.scenes {
//		if strings.Index(v.scenename, s.lastObjectName+":"+keyword) == 0 {
//			slice := strings.Split(v.scenename, ":")
//			if len(slice) < 3 {
//				slice = append(slice, "")
//			}
//			switchname := slice[2]
//			if switchname != "" && !GameInstance.gameSwitchManager.CheckSwitch(switchname) {
//				continue
//			}
//
//			s.PushScene()
//			//전체 스위치를 뒤져서 스위치도 조건에 넣어야함
//			s.ActiveScene(v)
//
//			return true
//		}
//	}
//
//	return false
//}

func (s *ScriptManager) GetSceneManager(scenename string) *SceneManager {
	v := s.FindScript(scenename)
	if v != nil {
		return v
	}

	v = &SceneManager{scenename: scenename}
	s.scenes = append(s.scenes, v)
	return v
}

func (s *ScriptManager) NewSceneManager(scenename string) *SceneManager {
	v := &SceneManager{scenename: scenename}
	s.scenes = append(s.scenes, v)
	return v
}

func (s *ScriptManager) TalkEnd() {
	s.lastActiveScene = nil
	s.lastObjectName = ""
}
