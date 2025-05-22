package main

import (
	"github.com/iamvkosarev/learning-cards/internal/app"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/server"
)

const cardsConfigPathEnvKey = "CARDS_CONFIG_PATH"

func main() {
	app.PrepareAndRunApp(
		server.BuildCardsServer, func() (config.Config, error) {
			return config.Load(cardsConfigPathEnvKey)
		},
	)
}
