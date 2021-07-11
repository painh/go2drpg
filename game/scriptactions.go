package game

import "image/color"

type ScriptActionInterface interface {
	Run()
}

type ScriptActionSetGameStatus struct {
	status int
}

func (s *ScriptActionSetGameStatus) Run() {
	GameInstance.SetStatus(s.status)
}

type ScriptActionText struct {
	text string
}

func (s *ScriptActionText) Run() {
	GameInstance.log.Add(color.RGBA{255, 255, 255, 255}, s.text)
}

type ScriptActionAddKeyword struct {
	keyword string
}

func (s *ScriptActionAddKeyword) Run() {
	GameInstance.keywordManager.Add(s.keyword)
}

type ScriptActionAddLocation struct {
	keyword string
}

func (s *ScriptActionAddLocation) Run() {
	GameInstance.player.ActiveLocation(s.keyword)
}

type ScriptActionAddPerson struct {
	keyword string
}

func (s *ScriptActionAddPerson) Run() {
	GameInstance.player.ActivePerson(s.keyword)
}


type ScriptActionSetSwitch struct {
	keyword string
	flag    bool
}

func (s *ScriptActionSetSwitch) Run() {
	GameInstance.gameSwitchManager.SetSwitch(s.keyword, s.flag)
}

type ScriptActionPlayMusic struct {
	filename string
}

func (s *ScriptActionPlayMusic) Run() {
	GameInstance.audio.PlayLoopMusic(s.filename)
}
