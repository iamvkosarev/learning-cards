package server

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/app"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/domain/contracts"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/client"
	server "github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc/interceptor/verification"
	sqlRepository "github.com/iamvkosarev/learning-cards/internal/infrastructure/repository/postgres"
	"github.com/iamvkosarev/learning-cards/internal/usecase"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type CardsClient interface {
	contracts.CardReader
	contracts.GroupReader
	Close()
}

type ReviewsServiceDeps struct {
	dbPool      *pgxpool.Pool
	cardsClient CardsClient
}
type ReviewsServer struct {
	ReviewsServiceDeps
}

func BuildReviewsServer(ctx context.Context, cfg *config.ReviewsConfig) (app.Server, error) {
	return BuildServer(
		ctx, cfg.Common,
		func(ctx context.Context, log *slog.Logger) (pb.ReviewServiceServer, ReviewsServiceDeps, error) {
			return prepareReviewService(ctx, cfg, log)
		}, pb.RegisterReviewServiceServer, pb.RegisterReviewServiceHandlerFromEndpoint, NewReviewsServer,
	)
}

func NewReviewsServer(serviceDeps ReviewsServiceDeps) *ReviewsServer {
	return &ReviewsServer{
		ReviewsServiceDeps: serviceDeps,
	}
}

func (a *ReviewsServer) Start() error {
	return nil
}

func (a *ReviewsServer) Shutdown(context.Context) {
	a.ReviewsServiceDeps.dbPool.Close()
	a.cardsClient.Close()
}

func prepareReviewService(ctx context.Context, cfg *config.ReviewsConfig, logger *slog.Logger) (
	pb.ReviewServiceServer,
	ReviewsServiceDeps,
	error,
) {
	dbPool, err := ConnectToDbPool(ctx, cfg.Database)
	if err != nil {
		return nil, ReviewsServiceDeps{}, err
	}

	progressRepo := sqlRepository.NewProgressRepository(dbPool)

	cardsClient, err := client.NewCardsClient(ctx, cfg.CardsServer.Address)
	if err != nil {
		return nil, ReviewsServiceDeps{}, err
	}

	verifier := contracts.VerifyFunc(verification.GetUserId)

	var reviewUseCase server.ReviewUseCase = usecase.NewReviewUseCase(
		usecase.ReviewUseCaseDeps{
			ProgressWriter: progressRepo,
			ProgressReader: progressRepo,
			CardReader:     cardsClient,
			GroupReader:    cardsClient,
		},
	)

	reviewsService := server.NewReviewService(
		server.ReviewServiceDeps{
			ReviewUseCase: reviewUseCase,
			AuthVerifier:  verifier,
			Logger:        logger,
		},
	)
	serviceDeps := ReviewsServiceDeps{dbPool: dbPool, cardsClient: cardsClient}
	return reviewsService, serviceDeps, nil
}
