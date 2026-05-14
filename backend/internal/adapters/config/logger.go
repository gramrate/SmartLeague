package config

import (
	"time"

	"github.com/spf13/viper"
)

type LoggerConfig interface {
	Debug() bool
	LogToFile() bool
	LogsDir() string
	TimeLocation() *time.Location
}

type LoggerConfigImpl struct {
	debug        bool
	logToFile    bool
	logsDir      string
	timeLocation *time.Location
}

func NewLoggerConfig() (*LoggerConfigImpl, error) {
	location, err := time.LoadLocation(viper.GetString("settings.timezone"))
	if err != nil {
		return nil, err
	}

	return &LoggerConfigImpl{
		debug:        viper.GetBool("settings.debug"),
		logToFile:    viper.GetBool("settings.logger.log-to-file"),
		logsDir:      viper.GetString("settings.logger.logs-dir"),
		timeLocation: location,
	}, nil
}

func (cfg *LoggerConfigImpl) Debug() bool {
	return cfg.debug
}

func (cfg *LoggerConfigImpl) LogToFile() bool {
	return cfg.logToFile
}

func (cfg *LoggerConfigImpl) LogsDir() string {
	return cfg.logsDir
}

func (cfg *LoggerConfigImpl) TimeLocation() *time.Location {
	return cfg.timeLocation
}
