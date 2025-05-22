package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/domain/contracts"
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

type CardsServer struct {
	config     config.Config
	grpcServer *grpc.Server
	httpServer *http.Server
	logger     *slog.Logger
	dbPool     *pgxpool.Pool
}

func BuildCardsServer(ctx context.Context, cfg config.Config) (*CardsServer, error) {
	return BuildServer(
		ctx, cfg, prepareCardService, pb.RegisterCardServiceServer, NewCardsServer,
	)
}

func NewCardsServer(
	cfg config.Config, grpcServer *grpc.Server, httpServer *http.Server,
	log *slog.Logger, dbPool *pgxpool.Pool,
) *CardsServer {
	return &CardsServer{
		config:     cfg,
		grpcServer: grpcServer,
		httpServer: httpServer,
		logger:     log,
		dbPool:     dbPool,
	}
}

func (a *CardsServer) Start() error {
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

func (a *CardsServer) Shutdown(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, a.config.Server.ShutdownTimeout)
	defer cancel()

	a.grpcServer.GracefulStop()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		a.logger.Error("HTTP server shutdown error", sl.Err(err))
	}

	a.dbPool.Close()
}

func prepareCardService(dbPool *pgxpool.Pool, logger *slog.Logger) (pb.CardServiceServer, error) {
	groupRepo := sqlRepository.NewGroupRepository(dbPool)
	cardRepo := sqlRepository.NewCardRepository(dbPool)

	verifier := contracts.VerifyFunc(verification.GetUserId)

	var groupUseCase server.GroupUseCase = usecase2.NewGroupUseCase(
		usecase2.GroupUseCaseDeps{
			GroupReader:  groupRepo,
			GroupWriter:  groupRepo,
			AuthVerifier: verifier,
		},
	)
	var cardsUseCase server.CardsUseCase = usecase2.NewCardsUseCase(
		usecase2.CardsUseCaseDeps{
			GroupReader:  groupRepo,
			CardWriter:   cardRepo,
			CardReader:   cardRepo,
			AuthVerifier: verifier,
		},
	)

	cardsService := server.NewCardService(
		server.CardServiceDeps{
			GroupUseCase: groupUseCase,
			CardsUseCase: cardsUseCase,
			Logger:       logger,
		},
	)
	return cardsService, nil
}
