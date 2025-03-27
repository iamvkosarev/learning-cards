package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Server struct {
	RestPrefix string `yaml:"rest_prefix"`
	RESTPort   string `yaml:"rest_port"`
	GRPCPort   string `yaml:"grpc_port"`
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

func MustLoad() *Config {
	path := os.Getenv("CONFIG_PATH")

	if path == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(path); err != nil {
		log.Fatalf("CONFIG_PATH does not exist at: %s\n", path)
	}
	var config Config
	err := cleanenv.ReadConfig(path, &config)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	return &config
}
