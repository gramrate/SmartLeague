package service_provider

import (
	"SmartLeague/internal/adapters/config"
	"fmt"
)

func (s *ServiceProvider) LoggerConfig() LoggerConfig {
	if s.loggerConfig == nil {
		cfg, err := config.NewLoggerConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get logger config: %w", err))
		}
		s.loggerConfig = cfg
	}

	return s.loggerConfig
}

func (s *ServiceProvider) PGConfig() PGConfig {
	if s.pgConfig == nil {
		var err error
		s.pgConfig, err = config.NewPGConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get PG config: %w", err))
		}
	}

	return s.pgConfig
}

func (s *ServiceProvider) RedisConfig() RedisConfig {
	if s.redisConfig == nil {
		var err error
		s.redisConfig, err = config.NewRedisConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get redis config: %w", err))
		}
	}

	return s.redisConfig
}

func (s *ServiceProvider) MinIOConfig() MinIOConfig {
	if s.minioConfig == nil {
		var err error
		s.minioConfig, err = config.NewMinIOConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get minIO config: %w", err))
		}
	}

	return s.minioConfig
}

func (s *ServiceProvider) ServerConfig() ServerConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get http config: %w", err))
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}
func (s *ServiceProvider) MailConfig() MailConfig {
	if s.mailConfig == nil {
		cfg, err := config.NewMailConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get mail config: %w", err))
		}

		s.mailConfig = cfg
	}

	return s.mailConfig
}

func (s *ServiceProvider) JWTConfig() JWTConfig {
	if s.jwtConfig == nil {
		cfg, err := config.NewJWTConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get JWT config: %w", err))
		}

		s.jwtConfig = cfg
	}

	return s.jwtConfig
}
