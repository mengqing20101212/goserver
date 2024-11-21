package gameServer

type EventType int

const (
	EventType_Login EventType = iota //玩家登录
	EventType_Logout
	EventType_CreateRole
)

type IEventListener interface {
	handleEvent(eventType EventType, params ...any)
	getRegisterEventType() []EventType
}

type EventProcess struct {
	eventListeners map[EventType][]IEventListener
}

func NewEventProcess() *EventProcess {
	return &EventProcess{
		eventListeners: make(map[EventType][]IEventListener),
	}
}
func (this *EventProcess) RegisterEvent(eventType EventType, listener IEventListener) {
	this.eventListeners[eventType] = append(this.eventListeners[eventType], listener)
}
func (this *EventProcess) HandleEvent(eventType EventType, params ...any) {
	listeners, ok := this.eventListeners[eventType]
	if !ok {
		return
	}
	for _, listener := range listeners {
		listener.handleEvent(eventType, params...)
	}
}
