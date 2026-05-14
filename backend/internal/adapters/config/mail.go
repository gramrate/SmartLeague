package config

import (
	"github.com/spf13/viper"
)

// ServerConfig defines an interface for HTTP server configuration.
type MailConfig interface {
	Host() string
	Port() int
	Mail() string
	Password() string
}

type MailConfigImpl struct {
	host     string
	port     int
	mail     string
	password string
}

// NewMailConfig initializes a new Mail configuration from environment variables.
func NewMailConfig() (*MailConfigImpl, error) {
	return &MailConfigImpl{
		host:     viper.GetString("service.mail.host"),
		port:     viper.GetInt("service.mail.port"),
		mail:     viper.GetString("service.mail.mail"),
		password: viper.GetString("service.mail.password"),
	}, nil
}

func (cfg *MailConfigImpl) Host() string {
	return cfg.host
}

func (cfg *MailConfigImpl) Port() int {
	return cfg.port
}

func (cfg *MailConfigImpl) Mail() string {
	return cfg.mail
}

func (cfg *MailConfigImpl) Password() string {
	return cfg.password
}
