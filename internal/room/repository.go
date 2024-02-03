package room

import (
	"anshulbansal02/scribbly/internal/repository"
	"anshulbansal02/scribbly/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RoomRepository struct {
	repository.Repository
	KeyMutex utils.KeyMutex
}

/********************** Repository Methods **********************/

func (r *RoomRepository) error(err error) error {
	return fmt.Errorf("room repository: %w", err)
}

// Create a new unsaved room
func (r *RoomRepository) NewRoom(adminId *string, roomType string) *Room {
	return &Room{
		ID:           generateRoomId(),
		Code:         generateRoomCode(),
		Participants: []string{*adminId},
		Type:         roomType,
		Admin:        adminId,
	}
}

// Save room to repository
func (r *RoomRepository) SaveRoom(ctx context.Context, room *Room) error {
	u, err := json.Marshal(room)
	if err != nil {
		return r.error(err)
	}

	if err = r.Rdb.Set(ctx, GetNamespaceKey(room.ID), u, 0).Err(); err != nil {
		return r.error(err)
	}

	return nil
}

// Get room by ID
func (r *RoomRepository) GetRoom(ctx context.Context, roomId string) (*Room, error) {
	k, err := r.Rdb.Get(ctx, GetNamespaceKey(roomId)).Result()
	if err != nil {
		return nil, r.error(ErrRoomNotFound)
	}

	room := &Room{}
	err = json.Unmarshal([]byte(k), room)
	if err != nil {
		return nil, r.error(err)
	}

	return room, nil
}

func (r *RoomRepository) RoomExists(ctx context.Context, roomId string) (bool, error) {
	exists, err := r.Rdb.Exists(ctx, GetNamespaceKey(roomId)).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, r.error(err)
	}

	return exists > 0, nil
}

func (r *RoomRepository) DeleteRoom(ctx context.Context, roomId string) error {
	err := r.Rdb.Del(ctx, GetNamespaceKey(roomId)).Err()
	if err != nil {
		return r.error(err)
	}

	return nil
}
