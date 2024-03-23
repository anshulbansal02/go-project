package chat

import (
	"anshulbansal02/scribbly/pkg/repository"
	"anshulbansal02/scribbly/pkg/utils"

	"context"
	"encoding/json"
	"fmt"
	"time"
)

type ChatRepository struct {
	repository.Repository
	KeyMutex utils.KeyMutex
}

/********************** Repository Methods **********************/

func (r *ChatRepository) error(err error) error {
	return fmt.Errorf("chat repository: %w", err)
}

func (r *ChatRepository) NewMessage(content string, userId string, conversationId string) *ChatMessage {
	return &ChatMessage{
		ID:             generateChatId(),
		Content:        content,
		UserId:         userId,
		Timestamp:      time.Now(),
		ConversationId: conversationId,
	}
}

func (r *ChatRepository) SaveMessage(ctx context.Context, msg *ChatMessage) error {
	m, err := json.Marshal(msg)
	if err != nil {
		return r.error(err)
	}

	err = r.Rdb.Set(ctx, getNamespaceKey(msg.ID), m, 0).Err()
	if err != nil {
		return r.error(err)
	}

	return nil
}
