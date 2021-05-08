package game

type FlowController struct {
	mainlooptoevent chan bool
	eventlooptomain chan bool
}

func (f *FlowController) Init() {
	f.mainlooptoevent = nil
	f.eventlooptomain = make(chan bool)
}

func (f *FlowController) StartEvent() {
	f.mainlooptoevent = make(chan bool)
}

func (f *FlowController) ShiftFlowToMainLoop() {
	f.eventlooptomain <- true
	//fmt.Println("main2event : wait")
	<-f.mainlooptoevent
	//fmt.Println("main2event : recv")
}

func (f *FlowController) ShiftFlowToEventLoop() bool {
	if f.mainlooptoevent == nil {
		return false
	}

	//fmt.Println("event2main : wait")
	_, ok := <-f.eventlooptomain
	//fmt.Println("event2main : arrived")
	f.mainlooptoevent <- true

	return ok
}

func (f *FlowController) EventEnd() {
	close(f.mainlooptoevent)
	f.mainlooptoevent = nil
	//fmt.Println("event end : send event2main")
	//f.eventlooptomain <- true
	//fmt.Println("event end : send event2main end")
}
