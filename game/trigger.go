package game

type Trigger interface {
	GetScripts() *[]ScriptAction
}

type CharacterCollisionTrigger struct {
}

func (s *CharacterCollisionTrigger) GetScripts() *[]ScriptAction {
	return nil
}

type TriggerManager struct {
}
