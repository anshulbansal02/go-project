package websockets

import "github.com/mitchellh/mapstructure"

type Event string

type UnpackedPayload map[string]any

func (p *UnpackedPayload) Assert(dst any) error {
	return mapstructure.Decode(p, dst)
}

type WebSocketMessage[T any] struct {
	Type      MessageType
	EventName Event
	Meta      map[string]any
	Payload   T
}

type OutgoingWebSocketMessage WebSocketMessage[any]

type IncomingWebSocketMessage WebSocketMessage[UnpackedPayload]

func NewNotification(name Event, payload any) OutgoingWebSocketMessage {
	return OutgoingWebSocketMessage{
		Type:      NotificationMessage,
		EventName: name,
		Payload:   payload,
	}
}

func NewResponse(name Event, requestId string, payload any) OutgoingWebSocketMessage {
	return OutgoingWebSocketMessage{
		Type:      ResponseMessage,
		EventName: name,
		Meta: map[string]any{
			"rId": requestId,
		},
		Payload: payload,
	}
}
