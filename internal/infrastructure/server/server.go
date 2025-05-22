package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/learning-cards/internal/app"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/database/postgres"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc/interceptor"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc/interceptor/verification"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/http/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
)

type NewServiceFunc[TService any, TServiceDeps any] func(context.Context, *slog.Logger) (TService, TServiceDeps, error)
type RegisterServiceFunc[TService any] func(grpc.ServiceRegistrar, TService)
type RegisterEndpoint func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
type NewServerFunc[TServer any, TServiceDeps any] func(TServiceDeps) TServer

type Wrapper struct {
	Server       app.Server
	serverConfig config.Server
	logger       *slog.Logger
	httpServer   *http.Server
	grpcServer   *grpc.Server
	verifier     verification.Verifier
}

func BuildServer[TService any, TServer app.Server, TServiceDeps any](
	ctx context.Context, cfg config.Config,
	newService NewServiceFunc[TService, TServiceDeps],
	registerService RegisterServiceFunc[TService],
	registerEndPoint RegisterEndpoint,
	newServer NewServerFunc[TServer, TServiceDeps],
) (
	*Wrapper,
	error,
) {
	logger, err := sl.SetupLogger(cfg.Env)
	if err != nil {
		return nil, fmt.Errorf("error setting up logger: %v", err)
	}

	service, serviceDeps, err := newService(ctx, logger)
	if err != nil {
		return nil, err
	}

	verifier, err := selectVerifier(cfg.SSO)
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

	err = registerEndPoint(
		ctx, gwMux, cfg.Server.GRPCPort,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to prepare HTTP gateway: %w", err)
	}

	httpMux := setupHTTPRouter(gwMux, cfg.Server)

	corsHandler := middleware.CorsWithOptions(httpMux, cfg.Server.CorsOptions)

	httpAddr := fmt.Sprintf("0.0.0.0%s", cfg.Server.RESTPort)

	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: corsHandler,
	}

	server := &Wrapper{
		Server:       newServer(serviceDeps),
		logger:       logger,
		serverConfig: cfg.Server,
		httpServer:   httpServer,
		grpcServer:   grpcServer,
		verifier:     verifier,
	}
	return server, nil
}

func (w *Wrapper) Start() error {
	go func() {
		w.Server.Start()
	}()

	lis, err := net.Listen("tcp", w.serverConfig.GRPCPort)
	if err != nil {
		return fmt.Errorf("error creating gRPC listener: %w", err)
	}

	go func() {
		w.logger.Info(fmt.Sprintf("Starting gRPC server on %s", w.serverConfig.GRPCPort))
		if err := w.grpcServer.Serve(lis); err != nil {
			w.logger.Error("failed to serve: %v", sl.Err(err))
		}
	}()

	w.logger.Info(
		fmt.Sprintf(
			"Starting REST gateway on %s%s/v%v/", w.httpServer.Addr, w.serverConfig.RestPrefix, w.serverConfig.Version,
		),
	)
	if err = w.httpServer.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to serve HTTP: %w", err)
		}
	}
	return nil
}

func (w *Wrapper) Shutdown(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, w.serverConfig.ShutdownTimeout)
	defer cancel()

	w.grpcServer.GracefulStop()
	w.verifier.Close()

	if err := w.httpServer.Shutdown(shutdownCtx); err != nil {
		w.logger.Error("HTTP server shutdown error", sl.Err(err))
	}

	w.Server.Shutdown(shutdownCtx)
}

func selectVerifier(ssoConfig config.SSO) (verification.Verifier, error) {
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

func setupHTTPRouter(gwMux *runtime.ServeMux, server config.Server) http.Handler {
	httpMux := http.NewServeMux()

	version := fmt.Sprintf("/v%v/", server.Version)
	httpMux.Handle(version, gwMux)

	httpMux.HandleFunc(
		server.RestPrefix+version, func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimPrefix(r.URL.Path, server.RestPrefix)
			r2 := new(http.Request)
			*r2 = *r
			r2.URL.Path = path
			gwMux.ServeHTTP(w, r2)
		},
	)

	return httpMux
}
