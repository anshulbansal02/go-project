package user

import (
	"anshulbansal02/scribbly/internal/repository"
	"anshulbansal02/scribbly/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type UserRepository struct {
	repository.Repository
	keyMutex utils.KeyMutex
}

func NewRepository(repo repository.Repository) *UserRepository {
	return &UserRepository{
		Repository: repo,
	}
}

/********************** Helper Functions **********************/

var generateUserId = utils.NewRandomStringGenerator(nil, 12)

func getNamespaceKey(userId string) string {
	return fmt.Sprintf("entity:user:%v", userId)
}

/********************** Repository Methods **********************/

// Create a new unsaved user
func (m *UserRepository) NewUser(name string) *User {
	return &User{ID: generateUserId(), Name: name}
}

// Save user to repository
func (m *UserRepository) SaveUser(ctx context.Context, user *User) error {

	u, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("user repository: %w", err)
	}

	err = m.Rdb.Set(ctx, getNamespaceKey(user.ID), u, 0).Err()
	if err != nil {
		return fmt.Errorf("user repository: %w", err)
	}

	return nil
}

// Get user by ID
func (m *UserRepository) GetUser(ctx context.Context, userId string) (*User, error) {
	u, err := m.Rdb.Get(ctx, getNamespaceKey(userId)).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("user respository: %w", err)
	}

	user := &User{}
	err = json.Unmarshal([]byte(u), user)
	if err != nil {
		return nil, fmt.Errorf("user respository: %w", err)
	}

	return user, nil
}

// Delete user by ID
func (m *UserRepository) DeleteUser(ctx context.Context, userId string) error {
	err := m.Rdb.Del(ctx, userId).Err()
	if err != nil {
		return fmt.Errorf("user repository: %w", err)
	}

	return nil
}

func (m *UserRepository) LockKey(userId string) func() {
	return m.keyMutex.Lock(userId)
}
