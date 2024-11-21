package gameServer

type HeroModule struct {
	BaseModule
}

func (self *HeroModule) InitModule(player *Player) {
	self.InitBaseModule(player)
}
func (this *HeroModule) GetModuleName() ModuleNameEnum {
	return ModuleNameHero
}

func (this *HeroModule) getRegisterEventType() []EventType {
	return []EventType{EventType_CreateRole, EventType_Login}
}

func (this *HeroModule) handleEvent(eventType EventType, params ...any) {
	switch eventType {
	case EventType_CreateRole:
		gameLogger.Info("HeroModule handleEvent CreateRole")
		break
	case EventType_Login:
		gameLogger.Info("HeroModule handleEvent Login")
		break
	}
}
