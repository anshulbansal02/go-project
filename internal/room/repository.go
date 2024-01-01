package room

import (
	"anshulbansal02/scribbly/internal/repository"
	"anshulbansal02/scribbly/pkg/utils"
	"context"
	"encoding/json"
)

type RoomRepository struct {
	repository.Repository
	keyMutex utils.KeyMutex
}

func NewRepository(repo repository.Repository) *RoomRepository {
	return &RoomRepository{
		Repository: repo,
	}
}

var generateRoomId = utils.NewRandomStringGenerator(nil, 8)
var generateRoomCode = utils.NewRandomStringGenerator(&utils.CHARSET_ALPHA_NUM, 6)

// Create a new unsaved room
func (m *RoomRepository) NewRoom(adminId *string, roomType string) *Room {
	return &Room{
		ID:    generateRoomId(),
		Code:  generateRoomCode(),
		Type:  roomType,
		Admin: adminId,
	}
}

// Save room to repository
func (m *RoomRepository) SaveRoom(ctx context.Context, room *Room) error {
	u, err := json.Marshal(room)
	if err != nil {
		return err
	}

	return m.Rdb.Set(ctx, GetNamespaceKey(room.ID), u, 0).Err()
}

// Get room by ID
func (m *RoomRepository) GetRoom(ctx context.Context, roomId string) (*Room, error) {
	r, err := m.Rdb.Get(ctx, GetNamespaceKey(roomId)).Result()
	if err != nil {
		return nil, ErrRoomNotFound
	}

	room := Room{}
	err = json.Unmarshal([]byte(r), &room)
	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (m *RoomRepository) RoomExists(ctx context.Context, roomId string) (bool, error) {
	exists, err := m.Rdb.Exists(ctx, GetNamespaceKey(roomId)).Result()
	return exists > 0, err
}

func (m *RoomRepository) DeleteRoom(ctx context.Context, roomId string) error {
	err := m.Rdb.Del(ctx, GetNamespaceKey(roomId)).Err()
	return err
}

// Lock roomId for mutex
func (m *RoomRepository) LockKey(roomId string) func() {
	return m.keyMutex.Lock(roomId)
}
