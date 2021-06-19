package game

type GameSwitchManager struct {
	dict map[string]bool
}

func (i *GameSwitchManager) Init() {
	i.dict = map[string]bool{}
}

func (i *GameSwitchManager) SetSwitch(key string, flag bool) {
	i.dict[key] = flag
}

func (i *GameSwitchManager) CheckSwitch(key string) bool {
	_, ok := i.dict[key]
	return ok
}
