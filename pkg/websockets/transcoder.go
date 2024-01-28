package websockets

import (
	"encoding/binary"
	"fmt"
	"math"

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

func EncodeMessage(message *OutgoingWebSocketMessage) ([]byte, error) {

	eventNameLength := len(message.EventName)
	if eventNameLength > maxEventNameLength {
		return nil, fmt.Errorf("websocket encode ~ event name length is greater than %d", maxEventNameLength)
	}

	meta, err := msgpack.Marshal(message.Meta)
	if err != nil {
		return nil, fmt.Errorf("websocket encode ~ failed to encode meta: %w", err)
	}

	metaLength := len(meta)
	if metaLength > math.MaxUint16 {
		return nil, fmt.Errorf("websocket encode ~ cannot encode meta data greater than 2 bytes in size")
	}

	payload, err := msgpack.Marshal(message.Payload)
	if err != nil {
		return nil, fmt.Errorf("websocket encode ~ failed to encode message payload: %w", err)
	}

	totalBytes := metaBytes + eventNameLength + metaLength + len(payload)
	encoded := make([]byte, totalBytes)

	// Encode Message Type
	encoded[0] = byte(message.Type) << 4

	// Encode Event Name Length
	encoded[1] = byte(eventNameLength)

	// Encode Event Name
	copy(encoded[2:], []byte(message.EventName))

	// Encode Meta Data Length
	offset := 2 + eventNameLength
	binary.LittleEndian.PutUint16(encoded[offset:2+offset], uint16(metaLength))

	// Encode Meta Data
	offset += 2
	copy(encoded[offset:], meta)

	// Encode Payload
	offset += len(meta)
	copy(encoded[offset:], payload)

	return encoded, nil
}

func DecodeMessage(encoded []byte) (*IncomingWebSocketMessage, error) {

	if len(encoded) < metaBytes {
		return nil, fmt.Errorf("websocket decode ~ message is too short to be decoded")
	}

	// Decode Message Type
	messageType := MessageType(encoded[0] >> 4)

	// Decode Event Name Length
	eventNameLength := int(encoded[1])

	// Decode Event Name
	offset := 2
	eventName := Event(string(encoded[offset : offset+eventNameLength]))

	// Decode Meta Data Length
	offset += int(eventNameLength)
	metaLength := int(binary.LittleEndian.Uint16(encoded[offset : 2+offset]))

	// Decode Meta Data
	offset += 2
	var meta map[string]any
	if err := msgpack.Unmarshal(encoded[offset:offset+metaLength], &meta); err != nil {
		return nil, err
	}

	// Decode Payload
	offset += metaLength
	var payload UnpackedPayload
	if err := msgpack.Unmarshal(encoded[offset:], &payload); err != nil {
		return nil, err
	}

	return &IncomingWebSocketMessage{
		Type:      messageType,
		EventName: eventName,
		Meta:      meta,
		Payload:   payload,
	}, nil
}
