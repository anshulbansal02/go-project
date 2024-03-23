package repository

import (
	"github.com/redis/go-redis/v9"
)

type Config struct {
	ServerAddress string
	Password      string
	DB            int
}

type Repository struct {
	Rdb *redis.Client
}

func New(config *Config) *Repository {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.ServerAddress,
		Password: config.Password,
		DB:       config.DB,
	})

	return &Repository{Rdb: redisClient}
}
