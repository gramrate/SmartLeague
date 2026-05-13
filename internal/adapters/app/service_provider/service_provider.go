package service_provider

import (
	"SmartLeague/internal/adapters/controller/api/validator"
	"SmartLeague/pkg/logger"
	"database/sql"

	"github.com/go-playground/form"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

type ServiceProvider struct {
	loggerConfig LoggerConfig
	pgConfig     PGConfig
	redisConfig  RedisConfig
	httpConfig   ServerConfig
	jwtConfig    JWTConfig
	mailConfig   MailConfig
	minioConfig  MinIOConfig

	redis *redis.Client
	minio *minio.Client
	sqlDB *sql.DB

	logger      *logger.Logger
	validator   *validator.Validator
	formDecoder *form.Decoder

	jwtService       jwtService
	tokenService     tokenService
	cookieService    cookieService
	userService      userService
	profileService   profileService
}

func New() *ServiceProvider {
	return &ServiceProvider{}
}
