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
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/auth"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/database/postgres"
	server "github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc/interceptor"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/http/middleware"
	sqlRepository "github.com/iamvkosarev/learning-cards/internal/infrastructure/repository/postgres"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	sso_pb "github.com/iamvkosarev/sso/pkg/proto/sso/v1"
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

	learningCardsServer, err := prepareLearningCardServer(cfg, dbPool, logger)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.SetupInterceptor(),
			interceptor.RecoveryInterceptor(logger),
			interceptor.LoggerUnaryServerInterceptor(logger),
			interceptor.ValidationInterceptor(logger),
		),
	)

	pb.RegisterLearningCardsServer(grpcServer, learningCardsServer)

	gwMux := runtime.NewServeMux()

	err = pb.RegisterLearningCardsHandlerFromEndpoint(
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

func (a *App) Run() error {
	lis, err := net.Listen("tcp", a.config.Server.GRPCPort)
	if err != nil {
		return fmt.Errorf("error creating gRPC listener: %w", err)
	}

	go func() {
		a.logger.Info(fmt.Sprintf("Starting gRPC server on %s", a.config.Server.GRPCPort))
		if err := a.grpcServer.Serve(lis); err != nil {
			a.logger.Error("failed to serve: %v", err)
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

func prepareLearningCardServer(cfg *config.Config, dbPool *pgxpool.Pool, logger *slog.Logger) (*server.Server, error) {
	groupRepo := sqlRepository.NewGroupRepository(dbPool)
	cardRepo := sqlRepository.NewCardRepository(dbPool)

	var authService contracts.AuthVerifier
	authService, err := getAuthVerifier(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to setup auth verifier: %w", err)
	}

	groupUseCase := usecase.NewGroupUseCase(
		usecase.GroupUseCaseDeps{
			GroupReader:  groupRepo,
			GroupWriter:  groupRepo,
			AuthVerifier: authService,
		},
	)
	cardsUseCase := usecase.NewCardsUseCase(
		usecase.CardsUseCaseDeps{
			GroupReader:  groupRepo,
			CardWriter:   cardRepo,
			CardReader:   cardRepo,
			AuthVerifier: authService,
		},
	)

	learningCardsServer := server.NewServer(groupUseCase, cardsUseCase, logger)
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

func getAuthVerifier(cfg *config.Config) (contracts.AuthVerifier, error) {
	var authService contracts.AuthVerifier
	if cfg.SSO.UseLocal {
		authService = auth.NewLocalService(cfg.SSO.LocalUserId)
		return authService, nil
	}

	ssoConn, err := grpc.NewClient(cfg.SSO.HostAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("error creating gRPC client: %w", err)
	}
	ssoClient := sso_pb.NewSSOClient(ssoConn)
	authService = auth.NewGRPCService(ssoClient)
	return authService, nil
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
