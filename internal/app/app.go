package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/learning-cards/internal/app/usecase"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/domain/contracts"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/database/postgres"
	server "github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc/interceptor"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc/interceptor/verification"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/http/middleware"
	sqlRepository "github.com/iamvkosarev/learning-cards/internal/infrastructure/repository/postgres"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
)

type App struct {
	config     *config.Config
	grpcServer *grpc.Server
	httpServer *http.Server
	logger     *slog.Logger
	dbPool     *pgxpool.Pool
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	logger, err := sl.SetupLogger(cfg.Env)
	if err != nil {
		return nil, fmt.Errorf("error setting up logger: %v", err)
	}

	dbPool, err := prepareDatabase(ctx, err)
	if err != nil {
		return nil, err
	}

	cardService, err := prepareCardService(dbPool, logger)
	if err != nil {
		return nil, err
	}

	verifier, err := selectVerifier(cfg)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.SetupInterceptor(),
			interceptor.RecoveryInterceptor(logger),
			verification.Interceptor(logger, verifier),
			interceptor.LoggerUnaryServerInterceptor(logger),
			interceptor.ValidationInterceptor(logger),
		),
	)

	pb.RegisterCardServiceServer(grpcServer, cardService)

	gwMux := runtime.NewServeMux()

	err = pb.RegisterCardServiceHandlerFromEndpoint(
		ctx, gwMux, cfg.Server.GRPCPort,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to prepare HTTP gateway: %w", err)
	}

	httpMux := setupHTTPRouter(gwMux, cfg.Server.RestPrefix)

	corsHandler := middleware.CorsWithOptions(httpMux, cfg.Server.CorsOptions)

	httpAddr := fmt.Sprintf("0.0.0.0%s", cfg.Server.RESTPort)

	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: corsHandler,
	}

	return &App{
		config:     cfg,
		grpcServer: grpcServer,
		httpServer: httpServer,
		logger:     logger,
		dbPool:     dbPool,
	}, nil
}

func selectVerifier(cfg *config.Config) (verification.Verifier, error) {
	if cfg.SSO.UseLocal {
		return verification.NewStubVerifier(cfg.SSO.LocalUserId), nil
	}
	return verification.NewGRPCVerifier(cfg.SSO.HostAddress)
}

func (a *App) Run() error {
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
	if err := a.httpServer.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to serve HTTP: %w", err)
		}
	}
	return nil
}

func (a *App) Shutdown(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, a.config.Server.ShutdownTimeout)
	defer cancel()

	a.grpcServer.GracefulStop()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		a.logger.Error("HTTP server shutdown error", sl.Err(err))
	}

	a.dbPool.Close()
}

func prepareCardService(dbPool *pgxpool.Pool, logger *slog.Logger) (*server.CardService, error) {
	groupRepo := sqlRepository.NewGroupRepository(dbPool)
	cardRepo := sqlRepository.NewCardRepository(dbPool)
	progressRepo := sqlRepository.NewProgressRepository(dbPool)

	verifier := contracts.VerifyFunc(verification.GetUserId)

	var groupUseCase server.GroupUseCase = usecase.NewGroupUseCase(
		usecase.GroupUseCaseDeps{
			GroupReader:  groupRepo,
			GroupWriter:  groupRepo,
			AuthVerifier: verifier,
		},
	)
	var cardsUseCase server.CardsUseCase = usecase.NewCardsUseCase(
		usecase.CardsUseCaseDeps{
			GroupReader:  groupRepo,
			CardWriter:   cardRepo,
			CardReader:   cardRepo,
			AuthVerifier: verifier,
		},
	)

	var progressUseCase server.ProgressUseCase = usecase.NewProgressUseCase(
		usecase.ProgressUseCaseDeps{
			ProgressReader: progressRepo,
		},
	)
	var reviewUseCase server.ReviewUseCase = usecase.NewReviewUseCase(
		usecase.ReviewUseCaseDeps{
			ProgressReader: progressRepo,
			ProgressWriter: progressRepo,
		},
	)

	learningCardsServer := server.NewServer(
		server.Deps{
			GroupUseCase:    groupUseCase,
			CardsUseCase:    cardsUseCase,
			ProgressUseCase: progressUseCase,
			ReviewUseCase:   reviewUseCase,
			Logger:          logger,
		},
	)
	return learningCardsServer, nil
}

func prepareDatabase(ctx context.Context, err error) (*pgxpool.Pool, error) {
	dns := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_SERVICE_NAME"),
		os.Getenv("DB_PORT_INTERNAL"), os.Getenv("DB_NAME"),
	)
	dbPool, err := postgres.NewPostgresPool(ctx, dns)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	return dbPool, nil
}

func setupHTTPRouter(gwMux *runtime.ServeMux, restPrefix string) http.Handler {
	httpMux := http.NewServeMux()

	const firstVersion = "/v1/"
	httpMux.Handle(firstVersion, gwMux)

	httpMux.HandleFunc(
		restPrefix+firstVersion, func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimPrefix(r.URL.Path, restPrefix)
			r2 := new(http.Request)
			*r2 = *r
			r2.URL.Path = path
			gwMux.ServeHTTP(w, r2)
		},
	)

	return httpMux
}
