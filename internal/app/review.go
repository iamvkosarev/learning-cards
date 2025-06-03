package app

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/client"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/server"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/server/interceptor/verification"
	"github.com/iamvkosarev/learning-cards/internal/repository/postgres"
	"github.com/iamvkosarev/learning-cards/internal/service"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type closer interface {
	Close()
}

type reviewsAppDeps struct {
	dbPool      *pgxpool.Pool
	cardsClient closer
}
type reviewsApp struct {
	reviewsAppDeps
}

func NewReviewsApp(ctx context.Context, cfg *config.ReviewsConfig) (*grpcAppWrapper, error) {
	return buildApp(
		ctx, cfg.Common,
		func(ctx context.Context, log *slog.Logger) (pb.ReviewServiceServer, reviewsAppDeps, error) {
			return prepareReviewServer(ctx, cfg, log)
		}, pb.RegisterReviewServiceServer, pb.RegisterReviewServiceHandlerFromEndpoint, newReviewsApp,
	)
}

func newReviewsApp(serviceDeps reviewsAppDeps) *reviewsApp {
	return &reviewsApp{
		reviewsAppDeps: serviceDeps,
	}
}

func (a *reviewsApp) start() error {
	return nil
}

func (a *reviewsApp) shutdown(context.Context) {
	a.reviewsAppDeps.dbPool.Close()
	a.cardsClient.Close()
}

func prepareReviewServer(ctx context.Context, cfg *config.ReviewsConfig, logger *slog.Logger) (
	pb.ReviewServiceServer,
	reviewsAppDeps,
	error,
) {
	dbPool, err := connectToDbPool(ctx, cfg.Database)
	if err != nil {
		return nil, reviewsAppDeps{}, err
	}

	progressRepo := postgres.NewProgressRepository(dbPool)

	cardsServerClient, err := client.NewCardsClient(ctx, cfg.CardsServer.Address)
	if err != nil {
		return nil, reviewsAppDeps{}, err
	}

	verifier := server.VerifyFunc(verification.GetUserId)

	var reviewService server.ReviewService = service.NewReviewService(
		service.ReviewServiceDeps{
			ProgressWriter: progressRepo,
			ProgressReader: progressRepo,
			CardReader:     cardsServerClient,
			GroupReader:    cardsServerClient,
			Config:         cfg.ReviewsService,
		},
	)

	server := server.NewReviewServer(
		server.ReviewServerDeps{
			ReviewService: reviewService,
			AuthVerifier:  verifier,
			Logger:        logger,
		},
	)
	appDeps := reviewsAppDeps{dbPool: dbPool, cardsClient: cardsServerClient}
	return server, appDeps, nil
}
