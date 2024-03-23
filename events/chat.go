package events

import "anshulbansal02/scribbly/pkg/websockets"

var Chat = struct {
	ChatMessage websockets.Event
}{
	ChatMessage: "chat_message",
}
