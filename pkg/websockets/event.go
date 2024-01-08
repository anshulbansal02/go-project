package websockets

type Event string

type WebSocketMessage struct {
	Type      MessageType
	EventName Event
	Payload   any
}

func NewEvent(name Event, payload any) WebSocketMessage {
	return WebSocketMessage{
		Type:      NotificationMessage,
		EventName: name,
		Payload:   payload,
	}
}
