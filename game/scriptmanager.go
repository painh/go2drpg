package game

//액션 종류
//대화
//선택
//선택에 따른 무언가
//변수 조절

const (
	SCRIPT_ACTION_TYPE_UNKNOWN = 0
	SCRIPT_ACTION_TYPE_TEXT    = 1
)

type ScriptAction interface {
	BreakUpdateLoop() bool
	Run()
	GetType() int
}

//func (s *ScriptAction) BreakUpdateLoop() bool {
//	return false
//}
//
//func (s *ScriptAction) Run() {
//}
//
//func (s *ScriptAction) GetType() int {
//	return SCRIPT_ACTION_TYPE_UNKNOWN
//}

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
	GameInstance.Log.AddWithPrompt(s.text)
}

type SceneManager struct {
	scene  []ScriptAction
	cursor int
	isOver bool
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
	s.cursor++

	currentScript.Run()

	return true

}

type ScriptManager struct {
	scripts      map[string]*SceneManager
	activeScript *SceneManager
}

func (s *ScriptManager) Init() {
	s.scripts = map[string]*SceneManager{}
	s.activeScript = nil

	scr := SceneManager{}
	scr.scene = append(scr.scene, &ScriptActionText{text: "안녕하세요."})
	s.Add("test", &scr)
}

func (s *ScriptManager) Update() {
	if s.activeScript == nil {
		return
	}

	s.activeScript.Update()
}

func (s *ScriptManager) Add(scriptName string, scr *SceneManager) {
	s.scripts[scriptName] = scr
}

func (s *ScriptManager) RunObjectScript(objName string) bool {
	objName = "test"
	val, exists := s.scripts[objName]
	if !exists {
		return false
	}

	s.activeScript = val
	s.activeScript.Start()

	return true
}
