package game

type Location struct {
	name     string
	location *LocationInfo
}

type Player struct {
	day        int
	curTimeMin int

	activeLocation map[string]Location
}

func (p *Player) Init() {
	p.activeLocation = make(map[string]Location)
	p.curTimeMin = SettingConfigInstance.StartTimeMin
}
func (p *Player) ActiveLocation(name string) {
	for i := 0; i < len(SettingConfigInstance.LocationList); i++ {
		v := &SettingConfigInstance.LocationList[i]
		if v.Name == name {
			p.activeLocation[name] = Location{name: name, location: v}
		}
	}
}

func (p *Player) AddTime(min int) {
	p.curTimeMin += min
}
