package chat

import (
	"anshulbansal02/scribbly/pkg/repository"
)

func NewRepository(repo repository.Repository) *ChatRepository {
	return &ChatRepository{
		Repository: repo,
	}
}

func NewService(chatRepository *ChatRepository) *ChatService {
	return &ChatService{
		chatRepository:     chatRepository,
		ChatMessageChannel: make(chan *ChatMessage),
	}
}

func SetupConcreteService(repository repository.Repository) *ChatService {
	chatService := NewService(
		NewRepository(repository),
	)

	return chatService
}
