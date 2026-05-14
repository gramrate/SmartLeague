package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type RedisConfig interface {
	Addr() string
	Password() string
	DB() int
}

type RedisConfigImpl struct {
	host     string
	port     int
	password string
	db       int
}

func NewRedisConfig() (*RedisConfigImpl, error) {
	return &RedisConfigImpl{
		host:     viper.GetString("service.redis.host"),
		port:     viper.GetInt("service.redis.port"),
		password: viper.GetString("service.redis.password"),
		db:       viper.GetInt("service.redis.db"),
	}, nil
}

func (cfg *RedisConfigImpl) Addr() string {
	return fmt.Sprintf("%s:%d", cfg.host, cfg.port)
}

func (cfg *RedisConfigImpl) Password() string {
	return cfg.password
}

func (cfg *RedisConfigImpl) DB() int {
	return cfg.db
}
