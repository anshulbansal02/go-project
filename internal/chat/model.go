package chat

import "time"

type ChatMessage struct {
	// Message Id
	ID string
	// Message content sent by the user
	Content string
	// Meta data associated with the user, to be used by pipes
	Meta map[string]any
	// Sender of the message
	UserId string
	// Received timestamp
	Timestamp time.Time
	// Message is part of which conversation or (room, group, channel, ...)
	ConversationId string
}
