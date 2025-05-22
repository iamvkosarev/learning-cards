package config

import (
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
	Version         int           `yaml:"version"`
	RESTPort        string        `yaml:"rest_port"`
	GRPCPort        string        `yaml:"grpc_port"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	CorsOptions     CorsOptions   `yaml:"cors"`
}

type SSO struct {
	HostAddress string `yaml:"host_address"`
	UseLocal    bool   `yaml:"use_local"`
	LocalUserId int64  `yaml:"local_user_id"`
}

type Database struct {
	ConnectionStringKey string `yaml:"connection_string_key"`
}

type Config struct {
	Env    string `yaml:"env" env-default:"development"`
	Server `yaml:"server"`
	SSO    `yaml:"sso"`
}

func Load[TConfig any](configEnvKey string) (*TConfig, error) {
	path := os.Getenv(configEnvKey)

	if path == "" {
		return nil, fmt.Errorf("%s environment variable not set", configEnvKey)
	}

	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("%s does not exist at: %s", configEnvKey, path)
	}
	var config TConfig
	err := cleanenv.ReadConfig(path, &config)
	if err != nil {
		return nil, fmt.Errorf("failed load config at path %s: %w", path, err)
	}
	return &config, nil
}
