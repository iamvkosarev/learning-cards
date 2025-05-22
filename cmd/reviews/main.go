package main

import (
	"github.com/iamvkosarev/learning-cards/internal/app"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/server"
)

const reviewsConfigPathEnvKey = "REVIEWS_CONFIG_PATH"

func main() {
	app.PrepareAndRunApp(
		server.BuildReviewsServer, func() (config.Config, error) {
			return config.Load(reviewsConfigPathEnvKey)
		},
	)
}
