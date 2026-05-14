package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type MinIOConfig interface {
	Endpoint() string
	AccessKey() string
	SecretKey() string
	BucketName() string
	SSL() bool
	Timeout() time.Duration
}

type MinIOConfigImpl struct {
	host       string
	port       int
	accessKey  string
	secretKey  string
	bucketName string
	ssl        bool
	timeout    time.Duration
}

func NewMinIOConfig() (*MinIOConfigImpl, error) {
	return &MinIOConfigImpl{
		host:       viper.GetString("service.minio.host"),
		port:       viper.GetInt("service.minio.port"),
		accessKey:  viper.GetString("service.minio.access-key"),
		secretKey:  viper.GetString("service.minio.secret-key"),
		bucketName: viper.GetString("service.minio.bucket-name"),
		ssl:        viper.GetBool("service.minio.ssl"),
		timeout:    viper.GetDuration("service.minio.timeout") * time.Second,
	}, nil
}

func (cfg *MinIOConfigImpl) Endpoint() string {
	return fmt.Sprintf("%s:%d", cfg.host, cfg.port)
}

func (cfg *MinIOConfigImpl) AccessKey() string {
	return cfg.accessKey
}

func (cfg *MinIOConfigImpl) SecretKey() string {
	return cfg.secretKey
}

func (cfg *MinIOConfigImpl) BucketName() string {
	return cfg.bucketName
}

func (cfg *MinIOConfigImpl) SSL() bool {
	return cfg.ssl
}

func (cfg *MinIOConfigImpl) Timeout() time.Duration {
	return cfg.timeout
}
