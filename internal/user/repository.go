package user

import (
	"anshulbansal02/scribbly/pkg/repository"
	"anshulbansal02/scribbly/pkg/utils"

	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type UserRepository struct {
	repository.Repository
	KeyMutex utils.KeyMutex
}

/********************** Repository Methods **********************/

func (r *UserRepository) error(err error) error {
	return fmt.Errorf("user repository: %w", err)
}

// Create a new unsaved user
func (r *UserRepository) NewUser(name string) *User {
	return &User{
		ID:   generateUserId(),
		Name: name,
	}
}

// Save user to repository
func (r *UserRepository) SaveUser(ctx context.Context, user *User) error {

	u, err := json.Marshal(user)
	if err != nil {
		return r.error(err)
	}

	err = r.Rdb.Set(ctx, getNamespaceKey(user.ID), u, 0).Err()
	if err != nil {
		return r.error(err)
	}

	return nil
}

func (r *UserRepository) GetUsers(ctx context.Context, userIds []string) ([]*User, error) {

	namespacedIds := utils.MapFunc[string, string](userIds, func(id string) string {
		return getNamespaceKey(id)
	})

	u, err := r.Rdb.MGet(ctx, namespacedIds...).Result()
	if err != nil {
		return nil, r.error(err)
	}

	users := []*User{}

	for _, encodedUser := range u {
		if encodedUser == nil {
			continue
		}
		user := &User{}
		err = json.Unmarshal([]byte(encodedUser.(string)), user)
		if err != nil {
			return nil, r.error(err)
		}

		users = append(users, user)
	}

	return users, nil
}

// Get a user by ID
func (r *UserRepository) GetUser(ctx context.Context, userId string) (*User, error) {
	u, err := r.Rdb.Get(ctx, getNamespaceKey(userId)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, r.error(ErrUserNotFound)
		}
		return nil, r.error(err)
	}

	user := &User{}
	err = json.Unmarshal([]byte(u), user)
	if err != nil {
		return nil, r.error(err)
	}

	return user, nil
}

// Delete user by ID
func (r *UserRepository) DeleteUser(ctx context.Context, userId string) error {
	err := r.Rdb.Del(ctx, getNamespaceKey(userId)).Err()
	if err != nil {
		return r.error(err)
	}

	return nil
}

func (r *UserRepository) UserExists(ctx context.Context, userId string) (bool, error) {
	exists, err := r.Rdb.Exists(ctx, getNamespaceKey(userId)).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, r.error(err)
	}

	return exists > 0, nil
}
