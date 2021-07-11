package game

import (
	"fmt"
	"image/color"
)

type Location struct {
	name     string
	location *LocationInfo
}

type Player struct {
	day        int
	curTimeMin int

	activeLocation      map[string]Location
	activePerson        map[string]bool
	currentLocationName string
}

func (p *Player) Init() {
	p.activeLocation = make(map[string]Location)
	p.activePerson = make(map[string]bool)
	p.curTimeMin = SettingConfigInstance.StartTimeMin
}

func (p *Player) ActiveLocation(name string) {
	_, ok := p.activeLocation[name]
	if ok {
		return
	}

	for i := 0; i < len(SettingConfigInstance.LocationList); i++ {
		v := &SettingConfigInstance.LocationList[i]
		if v.Name == name {
			p.activeLocation[name] = Location{name: name, location: v}
			GameInstance.log.AddWithPrompt(color.RGBA{0, 255, 0, 255}, name, "으로 이동 할 수 있게 되었습니다.")
		}
	}
}

func (p *Player) ActivePerson(name string) {
	_, ok := p.activePerson[name]
	if ok {
		return
	}

	p.activePerson[name] = true
	GameInstance.log.AddWithPrompt(color.RGBA{0, 255, 0, 255}, name, "(이)가 인물 사전에 추가되었습니다.")
}

func (p *Player) AddTime(min int) {
	for i := 0; i < min; i++ {
		p.curTimeMin++
		str := p.GetTimeString()
		// TODO : 이 루프 안에서 여러개의 스크립트가 실행되는 것을 대비해야함
		if GameInstance.scriptManager.FindScript(str) == nil {
			continue
		}
		GameInstance.scriptManager.RunObjectScript(str)

	}
	p.curTimeMin += min
}

func (p *Player) GetTimeString() string {
	hour := GameInstance.player.curTimeMin / 60
	min := GameInstance.player.curTimeMin % 60

	return fmt.Sprintf("%d/%02d:%02d", p.day, int(hour), int(min))
}
