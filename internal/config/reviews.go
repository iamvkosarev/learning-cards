package config

type CardsServer struct {
	Address string `json:"address"`
}

type ReviewsConfig struct {
	Common      Config `json:"common"`
	Database    `yaml:"database"`
	CardsServer `yaml:"cards_server"`
}
