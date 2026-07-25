package config

import (
	"context"
	"fmt"
	"os"

	"github.com/briheet/nozarashi/internal/specs"

	"github.com/go-playground/validator/v10"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func ParseTOMLConfig(ctx context.Context, filepath string) (*specs.Specs, error) {
	// Init a new viper path
	v := viper.New()

	// We only supporting toml for now
	v.SetConfigType("toml")

	// Read the exact configuration file passed by the command.
	v.SetConfigFile(filepath)
	v.AutomaticEnv()

	// Read config file from disk
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// Unmarshall config to struct
	var config specs.Specs
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Viper normalizes map keys, but container environment variables are case-sensitive.
	var rawConfig struct {
		Services map[string]struct {
			Environment map[string]string `toml:"environment"`
		} `toml:"services"`
	}

	configData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("read configuration file %q: %w", filepath, err)
	}
	if err := toml.Unmarshal(configData, &rawConfig); err != nil {
		return nil, fmt.Errorf("decode service environments: %w", err)
	}

	for serviceName, rawService := range rawConfig.Services {
		service := config.Services[serviceName]
		service.Environment = rawService.Environment
		config.Services[serviceName] = service
	}

	// Valdiate struct via validator package
	if err := validate.Struct(config); err != nil {
		return nil, err
	}

	return &config, nil
}
