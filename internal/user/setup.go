package user

import (
	"anshulbansal02/scribbly/internal/repository"
	tokenfactory "anshulbansal02/scribbly/pkg/token_factory"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Secret        []byte
	SigningMethod jwt.SigningMethod
}

func NewRepository(repo repository.Repository) *UserRepository {
	return &UserRepository{
		Repository: repo,
	}
}

func NewService(userRepo *UserRepository, tokenFactory *tokenfactory.TokenFactory[UserClaims]) *UserService {
	return &UserService{
		userRepo:     userRepo,
		tokenFactory: tokenFactory,
	}
}

func SetupConcreteService(repository repository.Repository, config Config) *UserService {
	userService := NewService(
		NewRepository(repository),
		tokenfactory.New[UserClaims](config.SigningMethod, config.Secret),
	)

	return userService
}
