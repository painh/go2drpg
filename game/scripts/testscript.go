package scripts

import (
	"fmt"
	"time"
)

var ScriptFunc = make(map[string]func())

type IEventProcessor interface {
	StartEvent()
	EndEvent()
	WaitForMainLoop()
	SetText(t string)
}

var eventProcessor IEventProcessor

//func waitmain() {
//	eventProcessor.WaitForMainLoop()
//}

func textwait(t string) {
	(eventProcessor).SetText(t)
	(eventProcessor).WaitForMainLoop()
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
		fmt.Println("start")
		textwait("헬로우")
		fmt.Println("1")
		time.Sleep(time.Second * 1)
		textwait("헬로우2222222")
		fmt.Println("2")
		time.Sleep(time.Second * 1)
		textwait("헬로우333333")
		fmt.Println("3")
	}
}
