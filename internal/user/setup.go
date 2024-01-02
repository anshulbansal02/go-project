package user

import (
	"anshulbansal02/scribbly/internal/repository"
	tokenfactory "anshulbansal02/scribbly/pkg/token_factory"
)

func NewRepository(repo repository.Repository) *UserRepository {
	return &UserRepository{
		Repository: repo,
	}
}

func NewService(userRepo *UserRepository, tokenFactory *tokenfactory.TokenFactory) *UserService {
	return &UserService{
		userRepo:     userRepo,
		tokenFactory: tokenFactory,
	}
}

func SetupConcreteService(repository repository.Repository, tokenFactory *tokenfactory.TokenFactory) *UserService {
	userService := NewService(NewRepository(repository), tokenFactory)

	return userService
}
