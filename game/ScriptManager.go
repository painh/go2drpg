package game

//액션 종류
//대화
//선택
//선택에 따른 무언가
//변수 조절

type ScriptAction interface {
	Update()
}

type ScriptActionText struct {
}

func (s *ScriptActionText) Update() {

}

type ScriptManager struct {
	list []*ScriptAction
}

func (s *ScriptManager) Init() {
}

func (s *ScriptManager) Update() {
}

func (s *ScriptManager) Add(scr *ScriptAction) {
	s.list = append(s.list, scr)
}
