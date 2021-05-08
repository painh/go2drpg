package scripts

import (
	"strconv"
)

var ScriptFunc = make(map[string]func())

type IEventProcessor interface {
	StartEvent()
	EndEvent()
	ShiftFlowToMainLoop()
	SetText(t string)
	TextSelect(t []string)
	GetLastSelectedIndex() int
	WaitOneFrameOn()
}

var eventProcessor IEventProcessor

//func waitmain() {
//	eventProcessor.ShiftFlowToMainLoop()
//}

func textwait(t string) {
	(eventProcessor).SetText(t)
	(eventProcessor).WaitOneFrameOn()
	(eventProcessor).ShiftFlowToMainLoop()
}

func textSelect(t []string) int {
	(eventProcessor).TextSelect(t)
	(eventProcessor).ShiftFlowToMainLoop()
	return (eventProcessor).GetLastSelectedIndex()
}

func StartEvent(eventname string) {
	go func() {
		(eventProcessor).StartEvent()
		ScriptFunc[eventname]()
		(eventProcessor).EndEvent()
	}()
}

func Init(ie IEventProcessor) {
	eventProcessor = ie
	ScriptFunc["slime"] = func() {
		textwait("헬로우")
		textwait("헬로우2222222")
		textwait("헬로우333333")
		i := textSelect(
			[]string{
				"1. 안녕", "2.하하하",
			})

		textwait("너의 선택 : " + strconv.Itoa(i))
	}
}
