package config

import (
	"errors"
	"os"
	"time"

	"github.com/alessandra1408/goqrlog/pkg/log"
	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Server struct {
	Port         string        `json:"port" validate:"required"`
	IdleTimeout  time.Duration `json:"idleTimeout" validate:"required"`
	ReadTimeout  time.Duration `json:"readTimeout" validate:"required"`
	WriteTimeout time.Duration `json:"writeTimeout" validate:"required"`
	SSLEnabled   bool          `json:"sslEnabled"`
	CertFile     string        `json:"certFile" validate:"required"`
	CertKey      string        `json:"certKey" validate:"required"`
}

type Auth struct {
	Key string `json:"key" validate:"required"`
}

type App struct {
	Name        string        `json:"name" validate:"required"`
	Environment string        `json:"environment" validate:"required"`
	HTTPTimeout time.Duration `json:"httpTimeout" validate:"required"`
}

type Database struct {
	Scheme         string `json:"scheme" validate:"required"`
	Host           string `json:"host" validate:"required"`
	Port           int    `json:"port" validate:"required"`
	SSLMode        string `json:"sslmode" validate:"required"`
	ChannelBinding string `json:"channelBinding" validate:"required"`
	Name           string `json:"databaseName" validate:"required"`
	Username       string `json:"username" validate:"required"`
	Password       string `json:"password" validate:"required"`
}

type Config struct {
	Auth     *Auth     `json:"auth" validate:"required"`
	App      *App      `json:"app" validate:"required"`
	Server   *Server   `json:"server" validate:"required"`
	Database *Database `json:"database" validate:"required"`
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
	v.SetDefault("app.httpTimeout", 30*time.Second)

	_ = v.BindEnv("app.name", "APP_NAME")
	_ = v.BindEnv("app.environment", "ENVIRONMENT")
	_ = v.BindEnv("auth.key", "AUTH_KEY")
	_ = v.BindEnv("server.idleTimeout", "IDLE_TIMEOUT")
	_ = v.BindEnv("server.readTimeout", "READ_TIMEOUT")
	_ = v.BindEnv("server.writeTimeout", "WRITE_TIMEOUT")
	_ = v.BindEnv("server.port", "PORT")
	_ = v.BindEnv("app.httpTimeout", "HTTP_TIMEOUT")
	_ = v.BindEnv("database.scheme", "DB_SCHEME")
	_ = v.BindEnv("database.host", "DB_HOST")
	_ = v.BindEnv("database.port", "DB_PORT")
	_ = v.BindEnv("database.username", "DB_USERNAME")
	_ = v.BindEnv("database.password", "DB_PASSWORD")
	_ = v.BindEnv("database.name", "DB_NAME")
	_ = v.BindEnv("database.sslmode", "DB_SSLMODE")
	_ = v.BindEnv("database.channelBinding", "DB_CHANNEL_BINDING")

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

func Log(log log.Log) error {
	ll, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		ll = "debug"
	}

	lv := log.Level()

	err := lv.Set(ll)
	if err != nil {
		return errors.New("invalid level: " + err.Error())
	}

	return nil
}

func Local() {
	l, ok := os.LookupEnv("LOCATION")
	if !ok {
		l = "America/Sao_Paulo"
	}

	time.Local, _ = time.LoadLocation(l)
}
