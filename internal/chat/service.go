package chat

import (
	"context"
)

type ChatService struct {
	chatRepository     *ChatRepository
	chatMessagePipes   []Pipe
	ChatMessageChannel chan *ChatMessage
}

type Pipe = func(ChatMessage) ChatMessage

func (s *ChatService) RegisterPipe(ctx context.Context, pipe Pipe) {
	s.chatMessagePipes = append(s.chatMessagePipes, pipe)
}

func (s *ChatService) CreateMessage(ctx context.Context, content, userId, conversationId string) (*ChatMessage, error) {

	msg := s.chatRepository.NewMessage(content, userId, conversationId)

	err := s.chatRepository.SaveMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
