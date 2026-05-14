package service_provider

import (
	"crypto/rsa"
	"time"
)

type LoggerConfig interface {
	Debug() bool
	LogToFile() bool
	LogsDir() string
	TimeLocation() *time.Location
}

type PGConfig interface {
	DSN() string
}

type RedisConfig interface {
	Addr() string
	Password() string
	DB() int
}

type ServerConfig interface {
	Address() string
	Port() int
	Host() string
	EnabledTLS() bool
	DevMode() bool
}

type JWTConfig interface {
	PrivateKey() *rsa.PrivateKey
	PublicKey() *rsa.PublicKey
	RefreshTokenExpires() time.Duration
	AccessTokenExpires() time.Duration
}

type MailConfig interface {
	Host() string
	Port() int
	Mail() string
	Password() string
}

type MinIOConfig interface {
	Endpoint() string
	AccessKey() string
	SecretKey() string
	BucketName() string
	SSL() bool
	Timeout() time.Duration
}
