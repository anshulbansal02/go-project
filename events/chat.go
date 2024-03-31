package events

import (
	"anshulbansal02/scribbly/pkg/websockets"
	"time"
)

var Chat = struct {
	ChatMessage websockets.Event
}{
	ChatMessage: "chat_message",
}

type IncomingChatMessageData struct {
	Content string         `msgpack:"content"`
	Meta    map[string]any `msgpack:"meta"`
}

type OutgoingChatMessageData struct {
	ID             int            `msgpack:"id"`
	Content        string         `msgpack:"content"`
	Meta           map[string]any `msgpack:"meta"`
	UserId         string         `msgpack:"userId"`
	Timestamp      time.Time      `msgpack:"timestamp"`
	ConversationId string         `msgpack:"conversationId"`
}
