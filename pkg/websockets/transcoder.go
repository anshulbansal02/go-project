package websockets

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
)

/*
Message Encoding

<Section:Length>

First two bytes are reserved for meta information

1. A message can be of 4 types and is set as first 4 bits of the encoded message.
 - Command
 - Notification
 - Request
 - Response
2. Length of the event name follows message type and is set as next 8 bits.
3. Next is the event name of length 2^x in which x is the value of last 8 bits.
4. Rest bits represent the payload/message which are encoded using msgpack

|<MessageType:4>|<Reserved:4>|<EventNameLength:8>|<EventName:2^0-2^8>|<Message>

*/

type MessageType int

const (
	NotificationMessage MessageType = 0x0
	CommandMessage      MessageType = 0x1
	RequestMessage      MessageType = 0x2
	ResponseMessage     MessageType = 0x3
)

func EncodeMessage(message *WebSocketMessage) ([]byte, error) {

	payload, err := msgpack.Marshal(message.Payload)
	if err != nil {
		return nil, fmt.Errorf("websocket encode ~ failed to encode message payload: %w", err)
	}

	eventNameLength := len(message.EventName)
	totalBits := metaBits + eventNameLength*8 + len(payload)*8
	totalBytes := (totalBits + 7) / 8 // Adds 7 to ceil result

	encoded := make([]byte, totalBytes)

	// Encode MessageType
	encoded[0] = byte(message.Type) << 4

	// Encode EventNameLength
	if eventNameLength > maxEventNameLength {
		return nil, fmt.Errorf("websocket encode ~ event name length is greater than %d", maxEventNameLength)
	}
	encoded[1] = byte(eventNameLength)

	// Encode EventName
	copy(encoded[2:], []byte(message.EventName))

	// Encode Payload
	copy(encoded[2+eventNameLength:], payload)

	return encoded, nil
}

func DecodeMessage(encoded []byte) (*WebSocketMessage, error) {

	if len(encoded) < 3 {
		return nil, fmt.Errorf("websocket decode ~ message is too short to be decoded")
	}

	// Decode MessageType
	messageType := MessageType(encoded[0] >> 4)

	// Decode EventNameLength
	eventNameLength := encoded[1]
	eventName := Event(string(encoded[2 : 2+eventNameLength]))

	// Decode Payload
	payload := encoded[2+eventNameLength:]

	return &WebSocketMessage{
		Type:      messageType,
		EventName: eventName,
		Payload:   payload,
	}, nil
}
