package user

import (
	tokenfactory "anshulbansal02/scribbly/pkg/token_factory"

	"context"
)

type UserService struct {
	userRepo     *UserRepository
	tokenFactory *tokenfactory.TokenFactory[UserClaims]
}

/********************** Service Methods **********************/

// Create a user using name
func (s *UserService) CreateUser(ctx context.Context, name string) (*User, error) {
	user := s.userRepo.NewUser(name)

	userSecret, err := s.tokenFactory.GenerateToken(UserClaims{
		UserId: user.ID,
	})
	if err != nil {
		return nil, err
	}
	user.Secret = userSecret

	err = s.userRepo.SaveUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUsers(ctx context.Context, userIds []string) ([]*User, error) {
	users, err := s.userRepo.GetUsers(ctx, userIds)
	return users, err
}

// Get a user by its Id
func (s *UserService) GetUser(ctx context.Context, userId string) (*User, error) {
	user, err := s.userRepo.GetUser(ctx, userId)
	return user, err
}

func (s *UserService) VerifyUserToken(token string) (*UserClaims, error) {

	claims := &UserClaims{}

	if err := s.tokenFactory.GetClaims(token, claims); err != nil {
		return nil, err
	}

	return claims, nil
}

// Delete a user by its Id
func (s *UserService) DeleteUser(ctx context.Context, userId string) error {
	err := s.userRepo.DeleteUser(ctx, userId)
	return err
}

// Upadte name of a user by its Id
func (s *UserService) UpdateUserName(ctx context.Context, userId string, newName string) error {
	defer s.userRepo.KeyMutex.Lock(userId)()

	user, err := s.userRepo.GetUser(ctx, userId)
	if err != nil {
		return err
	}

	user.Name = newName
	err = s.userRepo.SaveUser(ctx, user)

	return err
}

func (s *UserService) UserExists(ctx context.Context, userId string) (bool, error) {
	exists, err := s.userRepo.UserExists(ctx, userId)

	if err != nil {
		return false, err
	}

	return exists, nil
}
