package chat

import "time"

type ChatMessage struct {
	ID             int            `json:"id"`
	Content        string         `json:"content"`
	Meta           map[string]any `json:"meta"`
	UserId         string         `json:"userId"`
	Timestamp      time.Time      `json:"timestamp"`
	ConversationId string         `json:"conversationId"`
}
