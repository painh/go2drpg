package game

//const MOVE = "move"

type Cooldown struct {
	key        string
	activedsec int64
}

type CooldownManager struct {
	dict map[string]*Cooldown
}

func (c *CooldownManager) ActiveCooldown(cooldownkey string) {
	v, found := c.dict[cooldownkey]

	if !found {
		c.dict[cooldownkey] = &Cooldown{key: cooldownkey}
		v = c.dict[cooldownkey]
	}

	v.activedsec = makeTimestamp()
}

func (c *CooldownManager) IsCooldownOver(cooldownkey string, ms int64) bool {
	v, found := c.dict[cooldownkey]

	if !found {
		return true
	}

	now := makeTimestamp()

	if now-v.activedsec > ms {
		return true
	}

	return false
}
