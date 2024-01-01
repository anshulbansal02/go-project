package roomaggregates

import (
	"anshulbansal02/scribbly/internal/repository"
	"context"
)

type UserJoinRequestRepository struct {
	repository.Repository
}

func NewUserJoinRequestRepository(repository repository.Repository) *UserJoinRequestRepository {
	return &UserJoinRequestRepository{
		Repository: repository,
	}
}

const joinRequestsKey = "relation:user->room:join_requests"

func (r *UserJoinRequestRepository) CreateJoinRequest(ctx context.Context, userId string, roomId string) error {

	err := r.Rdb.HSet(ctx, joinRequestsKey, userId, roomId).Err()
	return err
}

func (r *UserJoinRequestRepository) GetUserJoinRequestedRoom(ctx context.Context, userId string) (string, error) {
	roomId, err := r.Rdb.HGet(ctx, joinRequestsKey, userId).Result()
	return roomId, err
}

func (r *UserJoinRequestRepository) DeleteJoinRequest(ctx context.Context, userId string) error {
	err := r.Rdb.HDel(ctx, joinRequestsKey, userId).Err()
	return err
}
