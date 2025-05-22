package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/database/postgres"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc/interceptor"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc/interceptor/verification"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/http/middleware"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

type NewServiceFunc[TService any] func(*pgxpool.Pool, *slog.Logger) (TService, error)
type RegisterServiceFunc[TService any] func(grpc.ServiceRegistrar, TService)
type NewServerFunc[TServer any] func(config.Config, *grpc.Server, *http.Server, *slog.Logger, *pgxpool.Pool) *TServer

func BuildServer[TService any, TServer any](
	ctx context.Context, cfg config.Config,
	newService NewServiceFunc[TService],
	registerService RegisterServiceFunc[TService],
	newServer NewServerFunc[TServer],
) (
	*TServer,
	error,
) {
	logger, err := sl.SetupLogger(cfg.Env)
	if err != nil {
		return nil, fmt.Errorf("error setting up logger: %v", err)
	}

	dbPool, err := ConnectToDbPool(ctx, cfg.Database)
	if err != nil {
		return nil, err
	}

	service, err := newService(dbPool, logger)
	if err != nil {
		return nil, err
	}

	verifier, err := SelectVerifier(cfg.SSO)
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

	registerService(grpcServer, service)

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

	return newServer(cfg, grpcServer, httpServer, logger, dbPool), nil
}

func SelectVerifier(ssoConfig config.SSO) (verification.Verifier, error) {
	if ssoConfig.UseLocal {
		return verification.NewStubVerifier(ssoConfig.LocalUserId), nil
	}
	return verification.NewGRPCVerifier(ssoConfig.HostAddress)
}

func ConnectToDbPool(ctx context.Context, database config.Database) (*pgxpool.Pool, error) {
	dns := os.Getenv(database.ConnectionStringKey)
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
