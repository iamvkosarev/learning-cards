package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/learning-cards/internal/app"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/domain/contracts"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/client"
	server "github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc/interceptor/verification"
	sqlRepository "github.com/iamvkosarev/learning-cards/internal/infrastructure/repository/postgres"
	usecase2 "github.com/iamvkosarev/learning-cards/internal/usecase"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"net/http"
)

type ReviewsApp struct {
	config     config.Config
	grpcServer *grpc.Server
	httpServer *http.Server
	logger     *slog.Logger
	dbPool     *pgxpool.Pool
}

func BuildReviewsServer(ctx context.Context, cfg config.Config) (app.Server, error) {
	return BuildServer(
		ctx, cfg, func(pool *pgxpool.Pool, log *slog.Logger) (pb.ReviewServiceServer, error) {
			return prepareReviewService(ctx, pool, log)
		}, pb.RegisterReviewServiceServer, NewReviewsServer,
	)
}

func NewReviewsServer(
	cfg config.Config, grpcServer *grpc.Server, httpServer *http.Server,
	log *slog.Logger, dbPool *pgxpool.Pool,
) *ReviewsApp {
	return &ReviewsApp{
		config:     cfg,
		grpcServer: grpcServer,
		httpServer: httpServer,
		logger:     log,
		dbPool:     dbPool,
	}
}

func (a *ReviewsApp) Start() error {
	lis, err := net.Listen("tcp", a.config.Server.GRPCPort)
	if err != nil {
		return fmt.Errorf("error creating gRPC listener: %w", err)
	}

	go func() {
		a.logger.Info(fmt.Sprintf("Starting gRPC server on %s", a.config.Server.GRPCPort))
		if err := a.grpcServer.Serve(lis); err != nil {
			a.logger.Error("failed to serve: %v", sl.Err(err))
		}
	}()

	a.logger.Info(fmt.Sprintf("Starting REST gateway on %s", a.httpServer.Addr))
	if err = a.httpServer.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to serve HTTP: %w", err)
		}
	}
	return nil
}

func (a *ReviewsApp) Shutdown(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, a.config.Server.ShutdownTimeout)
	defer cancel()

	a.grpcServer.GracefulStop()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		a.logger.Error("HTTP server shutdown error", sl.Err(err))
	}

	a.dbPool.Close()
}

func prepareReviewService(ctx context.Context, dbPool *pgxpool.Pool, logger *slog.Logger) (
	pb.ReviewServiceServer,
	error,
) {
	progressRepo := sqlRepository.NewProgressRepository(dbPool)

	cardsClient := client.NewCardsClient()

	verifier := contracts.VerifyFunc(verification.GetUserId)

	var reviewUseCase server.ReviewUseCase = usecase2.NewReviewUseCase(
		usecase2.ReviewUseCaseDeps{
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
	return reviewsService, nil
}
