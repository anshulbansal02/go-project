package room

import (
	"anshulbansal02/scribbly/internal/repository"
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const roomCodeIdMapKey = "map:room:code->id"

type RoomCodeIdMapRepository struct {
	repository.Repository
}

func (m *RoomCodeIdMapRepository) error(err error) error {
	return fmt.Errorf("room code->id map: %w", err)
}

func (m *RoomCodeIdMapRepository) GetRoomId(ctx context.Context, code string) (string, error) {
	roomId, err := m.Rdb.HGet(ctx, roomCodeIdMapKey, code).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", m.error(err)
	}

	return roomId, nil
}

func (m *RoomCodeIdMapRepository) Set(ctx context.Context, code string, roomId string) error {
	err := m.Rdb.HSet(ctx, roomCodeIdMapKey, code, roomId).Err()
	if err != nil {
		return m.error(err)
	}
	return nil
}

func (m *RoomCodeIdMapRepository) Delete(ctx context.Context, code string) error {
	err := m.Rdb.HDel(ctx, roomCodeIdMapKey, code).Err()
	if err != nil {
		return m.error(err)
	}
	return nil
}
