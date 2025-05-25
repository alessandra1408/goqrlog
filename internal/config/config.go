package config

import (
	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Auth struct {
	Key string `json:"key" validate:"required"`
}

type App struct {
	Environment string `json:"environment" validate:"required"`
}

type Config struct {
	Auth Auth `json:"auth" validate:"required"`
	App  App  `json:"app" validate:"required"`
}

func (cfg *Config) IsStaging() bool {
	return cfg.App.Environment == "staging"
}

func (cfg *Config) IsProduction() bool {
	return cfg.App.Environment == "production"
}

func Get() (*Config, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigType("json")
	v.SetConfigName("config")

	v.AutomaticEnv()

	_ = v.BindEnv("app.environment", "ENVIRONMENT")
	_ = v.BindEnv("auth.key", "AUTH_KEY")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	validate := validator.New()
	err = validate.Struct(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
