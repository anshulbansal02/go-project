package roomaggregates

import (
	"anshulbansal02/scribbly/pkg/repository"
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type UserJoinRequestRepository struct {
	repository.Repository
}

func (m *UserJoinRequestRepository) error(err error) error {
	return fmt.Errorf("room join requests repository: %w", err)
}

func (r *UserJoinRequestRepository) CreateJoinRequest(ctx context.Context, userId string, roomId string) error {
	err := r.Rdb.HSet(ctx, joinRequestsKey, userId, roomId).Err()
	if err != nil {
		return r.error(err)
	}
	return nil
}

func (r *UserJoinRequestRepository) GetUserJoinRequestedRoom(ctx context.Context, userId string) (string, error) {
	roomId, err := r.Rdb.HGet(ctx, joinRequestsKey, userId).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", r.error(repository.ErrEntityNotFound)
		}
		return "", r.error(err)
	}

	return roomId, nil
}

func (r *UserJoinRequestRepository) DeleteJoinRequest(ctx context.Context, userId string) error {
	err := r.Rdb.HDel(ctx, joinRequestsKey, userId).Err()
	if err != nil {
		return r.error(err)
	}
	return nil
}
