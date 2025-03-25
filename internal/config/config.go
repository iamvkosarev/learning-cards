package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Server struct {
	RESTPort string `yaml:"rest_port"`
	GRPCPort string `yaml:"grpc_port"`
}

type Config struct {
	Env    string `yaml:"env" env-default:"development"`
	Server `yaml:"server"`
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
