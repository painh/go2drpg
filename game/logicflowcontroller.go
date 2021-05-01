package game

import "fmt"

type FlowController struct {
	mainlooptoevent chan string
	eventlooptomain chan bool
}

func (f *FlowController) Init() {
	f.mainlooptoevent = nil
	f.eventlooptomain = make(chan bool)
}

func (f *FlowController) StartEvent() {
	f.mainlooptoevent = make(chan string)
}

func (f *FlowController) WaitForMainLoop() {
	f.eventlooptomain <- true
	fmt.Println("main2event : wait")
	<-f.mainlooptoevent
	fmt.Println("main2event : recv")
}

func (f *FlowController) WaitForEventLoop(op string) bool {
	if f.mainlooptoevent == nil {
		return false
	}

	fmt.Println("event2main : wait")
	_, ok := <-f.eventlooptomain
	fmt.Println("event2main : arrived")
	f.mainlooptoevent <- op

	return ok
}

func (f *FlowController) EventEnd() {
	close(f.mainlooptoevent)
	f.mainlooptoevent = nil
	fmt.Println("event end : send event2main")
	//f.eventlooptomain <- true
	fmt.Println("event end : send event2main end")
}
