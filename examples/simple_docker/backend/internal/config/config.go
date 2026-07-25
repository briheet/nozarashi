package config

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	API   APIConfig   `mapstructure:",squash" validate:"required"`
	DB    DBConfig    `mapstructure:",squash" validate:"required"`
	Redis RedisConfig `mapstructure:",squash" validate:"required"`
}

type APIConfig struct {
	Port              int    `mapstructure:"port" validate:"required"`
	CORSAllowedOrigin string `mapstructure:"cors_allowed_origin" validate:"required"`
	ReadHeaderTimeout int    `mapstructure:"read_header_timeout" validate:"required"`
	ReadTimeout       int    `mapstructure:"read_timeout" validate:"required"`
	WriteTimeout      int    `mapstructure:"write_timeout" validate:"required"`
	IdleTimeout       int    `mapstructure:"idle_timeout" validate:"required"`
}

type DBConfig struct {
	DatabaseURL     string        `mapstructure:"databaseurl" validate:"required"`
	MaxConns        int32         `mapstructure:"db_max_conns" validate:"required,min=1"`
	MinConns        int32         `mapstructure:"db_min_conns" validate:"required,min=1,ltefield=MaxConns"`
	MaxConnLifetime time.Duration `mapstructure:"db_max_conn_lifetime" validate:"required,gt=0"`
	MaxConnIdleTime time.Duration `mapstructure:"db_max_conn_idle_time" validate:"required,gt=0"`
}

type RedisConfig struct {
	Address  string `mapstructure:"redis_address" validate:"required"`
	Password string `mapstructure:"redis_password"`
	Database int    `mapstructure:"redis_database" validate:"min=0"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func LoadConfig(ctx context.Context, paths ...string) (*Config, error) {
	var err error
	var config Config
	if len(paths) == 0 || paths[0] == "" {
		return nil, fmt.Errorf("at least one config path is required")
	}

	// Viper config
	v := viper.New()
	v.SetConfigFile(paths[0])
	v.SetConfigType("env")

	// If we have already injected in the environment
	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	for _, path := range paths[1:] {
		if path == "" {
			return nil, fmt.Errorf("config path cannot be empty")
		}
		v.SetConfigFile(path)
		v.SetConfigType("env")
		if err := v.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("merge config %q: %w", path, err)
		}
	}

	err = v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	err = validate.Struct(config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
