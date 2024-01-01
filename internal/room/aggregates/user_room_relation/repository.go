package roomaggregates

import (
	"anshulbansal02/scribbly/internal/repository"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type UserRoomRelationRepository struct {
	repository.Repository
}

func NewUserRoomRelation(repository repository.Repository) *UserRoomRelationRepository {
	return &UserRoomRelationRepository{
		Repository: repository,
	}
}

func getUserToRoomRelationKey() string {
	return "relation:user->room"
}

func getRoomToUsersRelationKey(roomId string) string {
	return fmt.Sprintf("relation:room->users:%v", roomId)
}

func (r *UserRoomRelationRepository) GetRoomIdByUserId(ctx context.Context, userId string) (string, error) {
	roomId, err := r.Rdb.HGet(ctx, getUserToRoomRelationKey(), userId).Result()
	return roomId, err
}

func (r *UserRoomRelationRepository) GetUsersByRoomId(ctx context.Context, roomId string) ([]string, error) {
	users, err := r.Rdb.SMembers(ctx, getRoomToUsersRelationKey(roomId)).Result()
	return users, err
}

func (r *UserRoomRelationRepository) AddUserToRoom(ctx context.Context, roomId string, userId string) error {

	roomKey := getRoomToUsersRelationKey(roomId)
	userKey := getUserToRoomRelationKey()

	return r.Rdb.Watch(ctx, func(tx *redis.Tx) error {

		pipe := tx.Pipeline()

		pipe.SAdd(ctx, roomKey, userId)
		pipe.HSet(ctx, userKey, userId, roomId)

		_, err := pipe.Exec(ctx)

		return err

	}, userKey, roomKey)
}

func (r *UserRoomRelationRepository) RemoveUserFromRoom(ctx context.Context, roomId string, userId string) error {

	roomKey := getRoomToUsersRelationKey(roomId)
	userKey := getUserToRoomRelationKey()

	return r.Rdb.Watch(ctx, func(tx *redis.Tx) error {

		pipe := tx.Pipeline()

		pipe.SRem(ctx, roomKey, userId)
		pipe.HDel(ctx, userKey, userId, roomId)

		_, err := pipe.Exec(ctx)

		return err

	}, userKey, roomKey)

}
