package user

import (
	"context"
)

type UserService struct {
	userRepo *UserRepository
}

func NewService(userRepo *UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

/********************** Service Methods **********************/

// Create a user using name
func (s *UserService) CreateUser(ctx context.Context, name string) (*User, error) {
	user := s.userRepo.NewUser(name)

	err := s.userRepo.SaveUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Get a user by its Id
func (s *UserService) GetUser(ctx context.Context, userId string) (*User, error) {
	user, err := s.userRepo.GetUser(ctx, userId)
	return user, err
}

// Delete a user by its Id
func (s *UserService) DeleteUser(ctx context.Context, userId string) error {
	err := s.userRepo.DeleteUser(ctx, userId)
	return err
}

// Upadte name of a user by its Id
func (s *UserService) UpdateUserName(ctx context.Context, userId string, newName string) error {
	defer s.userRepo.LockKey(userId)()

	user, err := s.userRepo.GetUser(ctx, userId)
	if err != nil {
		return err
	}

	user.Name = newName
	err = s.userRepo.SaveUser(ctx, user)

	return err
}
