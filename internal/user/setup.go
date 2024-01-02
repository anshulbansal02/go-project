package user

import (
	"anshulbansal02/scribbly/internal/repository"
)

func NewRepository(repo repository.Repository) *UserRepository {
	return &UserRepository{
		Repository: repo,
	}
}

func NewService(userRepo *UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func SetupConcreteService(repository repository.Repository) *UserService {
	userService := NewService(NewRepository(repository))

	return userService
}
