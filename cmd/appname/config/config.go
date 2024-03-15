package config

import (
	"github.com/spf13/viper"

	"github.com/idan-fishman/go-echo-starter/pkg/validation"
)

type Log struct {
	Level string `validate:"required,oneof=debug info warn error panic fatal"`
}

type Redis struct {
	Database int    `validate:"required,min=0,max=15"`
	Host     string `validate:"required,hostname"`
	Password string `validate:"omitempty"`
	Port     uint16 `validate:"required,min=1,max=65535"`
}

type Server struct {
	GracefulShutdownTimeoutSeconds uint16 `validate:"required,min=0,max=60"`
	Port                           uint16 `validate:"required,min=1,max=65535"`
	RequestTimeoutSeconds          uint16 `validate:"required,min=1,max=60"`
}

type Config struct {
	Log    Log    `validate:"required"`
	Server Server `validate:"required"`
	Redis  Redis  `validate:"required"`
}

// Load reads the configuration from the environment variables and validates it.
func Load(v *validation.Validator) (Config, error) {
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("SERVER_GRACEFUL_SHUTDOWN_TIMEOUT_SEC", 10)
	viper.SetDefault("SERVER_PORT", 1312)
	viper.SetDefault("SERVER_REQUEST_TIMEOUT_SEC", 30)
	viper.SetDefault("REDIS_DATABASE", 0)
	viper.SetDefault("REDIS_PORT", 6379)

	// Load the configuration from the environment variables.
	config := Config{
		Log: Log{
			Level: viper.GetString("LOG_LEVEL"),
		},
		Server: Server{
			GracefulShutdownTimeoutSeconds: viper.GetUint16("SERVER_GRACEFUL_SHUTDOWN_TIMEOUT_SEC"),
			Port:                           viper.GetUint16("SERVER_PORT"),
			RequestTimeoutSeconds:          viper.GetUint16("SERVER_REQUEST_TIMEOUT_SEC"),
		},
		Redis: Redis{
			Database: viper.GetInt("REDIS_DATABASE"),
			Host:     viper.GetString("REDIS_HOST"),
			Password: viper.GetString("REDIS_PASSWORD"),
			Port:     viper.GetUint16("REDIS_PORT"),
		},
	}

	// Ensure the configuration is valid.
	err := v.V.Struct(config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
