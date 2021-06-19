package game

import "strings"

//액션 종류
//대화
//선택
//선택에 따른 무언가
//변수 조절

const (
	SCRIPT_ACTION_TYPE_UNKNOWN               = 0
	SCRIPT_ACTION_TYPE_TEXT                  = 1
	SCRIPT_ACTION_TYPE_SET_GAME_STATUS       = 2
	SCRIPT_ACTION_TYPE_ADD_KEYWORD_TO_PLAYER = 3
	SCRIPT_ACTION_TYPE_SET_SWITCH            = 5
)

type ScriptActionInterface interface {
	BreakUpdateLoop() bool
	Run()
	GetType() int
}

type ScriptActionSetGameStatus struct {
	status int
}

func (s *ScriptActionSetGameStatus) GetType() int {
	return SCRIPT_ACTION_TYPE_SET_GAME_STATUS
}

func (s *ScriptActionSetGameStatus) BreakUpdateLoop() bool {
	return false
}

func (s *ScriptActionSetGameStatus) Run() {
	GameInstance.SetStatus(s.status)
}

type ScriptActionText struct {
	text string
}

func (s *ScriptActionText) GetType() int {
	return SCRIPT_ACTION_TYPE_TEXT
}

func (s *ScriptActionText) BreakUpdateLoop() bool {
	return false
}

func (s *ScriptActionText) Run() {
	GameInstance.log.Add(s.text)
}

type ScriptActionKeyword struct {
	keyword string
}

func (s *ScriptActionKeyword) GetType() int {
	return SCRIPT_ACTION_TYPE_ADD_KEYWORD_TO_PLAYER
}

func (s *ScriptActionKeyword) BreakUpdateLoop() bool {
	return false
}

func (s *ScriptActionKeyword) Run() {
	GameInstance.keywordManager.Add(s.keyword)
}

type ScriptActionSetSwitch struct {
	keyword string
	flag    bool
}

func (s *ScriptActionSetSwitch) GetType() int {
	return SCRIPT_ACTION_TYPE_SET_SWITCH
}

func (s *ScriptActionSetSwitch) BreakUpdateLoop() bool {
	return false
}

func (s *ScriptActionSetSwitch) Run() {
	GameInstance.gameSwitchManager.SetSwitch(s.keyword, s.flag)
}

type SceneManager struct {
	scenename string
	scene     []ScriptActionInterface
	cursor    int
	isOver    bool
	person    bool
}

func (s *SceneManager) Start() {
	s.isOver = false
	s.cursor = 0
}

func (s *SceneManager) Update() bool {
	if s.isOver {
		return false
	}

	currentScript := s.scene[s.cursor]
	currentScript.Run()

	s.cursor++
	if s.cursor >= len(s.scene) {
		s.isOver = true
	}

	return true

}

type ScriptManager struct {
	scripts                map[string]*SceneManager
	activeScript           *SceneManager
	activeScriptStack      []*SceneManager
	lastObjectName         string
	invalidKeywordResponse string
}

func (s *ScriptManager) GetInvalidKeywordResponse() string {
	if s.invalidKeywordResponse == "" {
		return "무슨 말인지 모르겠군요."
	}

	return s.invalidKeywordResponse
}

func (s *ScriptManager) Init() {
	s.scripts = map[string]*SceneManager{}
	s.activeScriptStack = []*SceneManager{}
	s.activeScript = nil
	s.invalidKeywordResponse = ""
	//
	//// Test Code
	//scr := &SceneManager{person: true}
	//scr.scene = append(scr.scene, &ScriptActionText{text: "안녕하세요."})
	//s.Add("test", scr)
	//
	//scr.scene = append(scr.scene, &ScriptActionSetGameStatus{status: 2})
	//s.Add("test", scr)
	//
	//scr.scene = append(scr.scene, &ScriptActionKeyword{keyword: "살해도구"})
	//s.Add("test", scr)
	//
	//scr.scene = append(scr.scene, &ScriptActionKeyWordReaction{keyword: "살해도구"})
	//s.Add("test", scr)
	//
	//scr = &SceneManager{person: true}
	//scr.scene = append(scr.scene, &ScriptActionText{text: "난 그런거 모르는데?222222"})
	//s.Add("test:살해도구:testflag", scr)
	//
	//scr = &SceneManager{person: true}
	//scr.scene = append(scr.scene, &ScriptActionText{text: "난 그런거 모르는데?"})
	//s.Add("test:살해도구:", scr)
	//
	//scr.scene = append(scr.scene, &ScriptActionSetSwitch{"testflag", true})
	//s.Add("test:살해도구:", scr)

}

func (s *ScriptManager) Update() {
	if s.activeScript == nil {
		return
	}

	if s.activeScript.Update() == false {
		s.activeScript = s.PopScene()
	}
}

//
//func (s *ScriptManager) Add(scenename string, scr *SceneManager) {
//	s.scripts[scenename] = scr
//	s.scripts[scenename].scenename = scenename
//}

func (s *ScriptManager) RunObjectScript(objName string) bool {
	s.invalidKeywordResponse = ""

	val, exists := s.scripts[objName]
	if !exists {
		return false
	}

	s.lastObjectName = objName
	s.ActiveScene(val)

	if s.activeScript.person {
		GameInstance.log.AddWithPrompt(objName, "(와)과 대화를 시작합니다.")
	}

	return true
}

func (s *ScriptManager) ActiveScene(scene *SceneManager) bool {
	s.activeScript = scene
	s.activeScript.Start()

	return true
}

func (s *ScriptManager) PushScene() {
	if s.activeScript == nil {
		return
	}
	s.activeScriptStack = append(s.activeScriptStack, s.activeScript)
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

func (s *ScriptManager) RunKeywordScript(keyword string) bool {
	for k, v := range s.scripts {
		if strings.Index(k, s.lastObjectName+":"+keyword) == 0 {
			slice := strings.Split(k, ":")
			if len(slice) < 3 {
				slice = append(slice, "")
			}
			switchname := slice[2]
			if switchname != "" && !GameInstance.gameSwitchManager.CheckSwitch(switchname) {
				continue
			}

			s.PushScene()
			//전체 스위치를 뒤져서 스위치도 조건에 넣어야함
			s.ActiveScene(v)

			return true
		}
	}

	return false
}

func (s *ScriptManager) GetSceneManager(scenename string) *SceneManager {
	v, ok := s.scripts[scenename]
	if ok {
		return v
	}

	s.scripts[scenename] = &SceneManager{scenename: scenename}
	return s.scripts[scenename]
}
