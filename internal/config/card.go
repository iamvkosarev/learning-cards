package config

type CardsConfig struct {
	Common   Config `json:"common"`
	Database `yaml:"database"`
}
