package user

import (
	"anshulbansal02/scribbly/internal/repository"
)

func SetupConcreteService(repository repository.Repository) *UserService {
	userService := NewService(NewRepository(repository))

	return userService
}
