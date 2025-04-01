package config

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type CorsOptions struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

type Server struct {
	RestPrefix      string        `yaml:"rest_prefix"`
	RESTPort        string        `yaml:"rest_port"`
	GRPCPort        string        `yaml:"grpc_port"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	CorsOptions     CorsOptions   `yaml:"cors"`
}

type SSO struct {
	HostAddress string `yaml:"host_address"`
	UseLocal    bool   `yaml:"use_local"`
	LocalUserId int    `yaml:"local_user_id"`
}

type Config struct {
	Env    string `yaml:"env" env-default:"development"`
	Server `yaml:"server"`
	SSO    `yaml:"sso"`
}

func Load() (*Config, error) {
	path := os.Getenv("CONFIG_PATH")

	if path == "" {

		return nil, errors.New("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("CONFIG_PATH does not exist at: %s", path)
	}
	var config Config
	err := cleanenv.ReadConfig(path, &config)
	if err != nil {
		return nil, fmt.Errorf("failed load config at path %s: %w", path, err)
	}
	return &config, nil
}
