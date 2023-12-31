package config

type Config struct {
	AppServerPort int

	RedisServerAddress  string
	RedisServerPassword string
	RedisServerDBNumber int

	SigningKey string
}
