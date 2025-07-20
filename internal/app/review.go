package app

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/client"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/server"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/server/interceptor/verification"
	"github.com/iamvkosarev/learning-cards/internal/module"
	"github.com/iamvkosarev/learning-cards/internal/otel/tracing"
	"github.com/iamvkosarev/learning-cards/internal/repository/postgres"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/sdk/trace"
	"log/slog"
)

type closer interface {
	Close()
}

type reviewsAppDeps struct {
	dbPool        *pgxpool.Pool
	cardsClient   closer
	traceProvider *trace.TracerProvider
}
type reviewsApp struct {
	*reviewsAppDeps
}

func NewReviewsApp(ctx context.Context, cfg *config.ReviewsConfig) (*grpcAppWrapper, error) {
	return buildApp(
		ctx, cfg.Common,
		func(ctx context.Context, log *slog.Logger) (pb.ReviewServiceServer, *reviewsAppDeps, error) {
			return prepareReviewServer(ctx, cfg, log)
		}, pb.RegisterReviewServiceServer, pb.RegisterReviewServiceHandlerFromEndpoint, newReviewsApp,
	)
}

func newReviewsApp(serviceDeps *reviewsAppDeps) *reviewsApp {
	return &reviewsApp{
		reviewsAppDeps: serviceDeps,
	}
}

func (a *reviewsApp) start() error {
	return nil
}

func (a *reviewsApp) shutdown(ctx context.Context) error {
	a.reviewsAppDeps.dbPool.Close()
	a.cardsClient.Close()
	return a.traceProvider.Shutdown(ctx)
}

func prepareReviewServer(ctx context.Context, cfg *config.ReviewsConfig, logger *slog.Logger) (
	pb.ReviewServiceServer,
	*reviewsAppDeps,
	error,
) {
	traceProvider, err := tracing.SetupTracingProvider(ctx, cfg.Tracing)
	if err != nil {
		return nil, nil, err
	}

	dbPool, err := connectToDbPool(ctx, cfg.Database)
	if err != nil {
		return nil, nil, err
	}

	reviewRepo := postgres.NewReviewRepository(dbPool)

	cardsServerClient, err := client.NewCardsClient(ctx, cfg.CardsServer.Address)
	if err != nil {
		return nil, nil, err
	}

	verifier := server.VerifyFunc(verification.GetUserId)

	var reviewService = module.NewReviews(
		module.ReviewsDeps{
			ReviewWriter: reviewRepo,
			ReviewReader: reviewRepo,
			UserVerifier: verifier,
			CardReader:   cardsServerClient,
			GroupReader:  cardsServerClient,
			Config:       cfg.ReviewsService,
		},
	)

	server := server.NewReviewServer(
		server.ReviewServerDeps{
			ReviewService: reviewService,
			Logger:        logger,
		},
	)
	appDeps := &reviewsAppDeps{dbPool: dbPool, cardsClient: cardsServerClient, traceProvider: traceProvider}
	return server, appDeps, nil
}
