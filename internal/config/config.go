package config

import (
	"context"

	"github.com/briheet/nozarashi/internal/specs"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func ParseTOMLConfig(ctx context.Context, filepath string) (*specs.Specs, error) {
	// Init a new viper path
	v := viper.New()

	// We only supporting toml for now
	v.SetConfigType("toml")

	// Set the config file path
	v.AddConfigPath(filepath)
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

	// Valdiate struct via validator package
	if err := validate.Struct(config); err != nil {
		return nil, err
	}

	return &config, nil
}
