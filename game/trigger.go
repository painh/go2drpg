package game

type Trigger interface {
	GetScripts() *[]ScriptActionInterface
}

type CharacterCollisionTrigger struct {
}

func (s *CharacterCollisionTrigger) GetScripts() *[]ScriptActionInterface {
	return nil
}

type TriggerManager struct {
}
