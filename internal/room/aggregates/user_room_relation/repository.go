package roomaggregates

import (
	"anshulbansal02/scribbly/internal/repository"
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type UserRoomRelationRepository struct {
	repository.Repository
}

func (m *UserRoomRelationRepository) error(err error) error {
	return fmt.Errorf("user room relation repository: %w", err)
}

func (r *UserRoomRelationRepository) GetRoomIdByUserId(ctx context.Context, userId string) (string, error) {
	roomId, err := r.Rdb.HGet(ctx, getUserToRoomRelationKey(), userId).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", r.error(err)
	}
	return roomId, nil
}

func (r *UserRoomRelationRepository) GetUsersByRoomId(ctx context.Context, roomId string) ([]string, error) {
	users, err := r.Rdb.SMembers(ctx, getRoomToUsersRelationKey(roomId)).Result()
	if err != nil {
		return []string{}, r.error(err)
	}
	return users, nil
}

func (r *UserRoomRelationRepository) AddUserToRoom(ctx context.Context, roomId string, userId string) error {

	roomKey := getRoomToUsersRelationKey(roomId)
	userKey := getUserToRoomRelationKey()

	err := r.Rdb.Watch(ctx, func(tx *redis.Tx) error {

		pipe := tx.Pipeline()

		pipe.SAdd(ctx, roomKey, userId)
		pipe.HSet(ctx, userKey, userId, roomId)

		_, err := pipe.Exec(ctx)

		return err

	}, userKey, roomKey)

	if err != nil {
		return r.error(err)
	}
	return nil
}

func (r *UserRoomRelationRepository) RemoveUserFromRoom(ctx context.Context, roomId string, userId string) error {

	roomKey := getRoomToUsersRelationKey(roomId)
	userKey := getUserToRoomRelationKey()

	err := r.Rdb.Watch(ctx, func(tx *redis.Tx) error {

		pipe := tx.Pipeline()

		pipe.SRem(ctx, roomKey, userId)
		pipe.HDel(ctx, userKey, userId, roomId)

		_, err := pipe.Exec(ctx)

		return err

	}, userKey, roomKey)

	if err != nil {
		return r.error(err)
	}
	return nil
}
