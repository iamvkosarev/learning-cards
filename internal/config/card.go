package config

import (
	"time"
)

type CardsConfig struct {
	Common          Config `yaml:"common"`
	Database        `yaml:"database"`
}
