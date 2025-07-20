package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/http/middleware"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/server/interceptor"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/server/interceptor/verification"
	"github.com/iamvkosarev/learning-cards/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type app interface {
	start() error
	shutdown(context.Context) error
}

type newServerFunc[TServer any, TAppDeps any] func(context.Context, *slog.Logger) (TServer, TAppDeps, error)
type registerServerFunc[TServer any] func(grpc.ServiceRegistrar, TServer)
type registerEndpoint func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
type newAppFunc[TApp any, TAppDeps any] func(TAppDeps) TApp

type grpcAppWrapper struct {
	app          app
	serverConfig config.Server
	logger       *slog.Logger
	httpServer   *http.Server
	grpcServer   *grpc.Server
	verifier     verification.Verifier
}

func (w *grpcAppWrapper) Run(
	ctx context.Context,
) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := w.start(); err != nil {
			log.Printf("Failed to initialize server: %v", err)
			shutdown <- syscall.SIGTERM
		}
	}()

	<-shutdown
	log.Println("Shutting down gracefully")

	err := w.shutdown(ctx)
	if err != nil {
		log.Printf("Failed to gracefully shutdown server: %v", err)
	}
	log.Println("app stopped")
}

func buildApp[TServer any, TApp app, TAppDeps any](
	ctx context.Context, cfg config.Config,
	newServerFunc newServerFunc[TServer, TAppDeps],
	registerServerFunc registerServerFunc[TServer],
	registerEndPointFunc registerEndpoint,
	newAppFunc newAppFunc[TApp, TAppDeps],
) (
	*grpcAppWrapper,
	error,
) {
	logger, err := sl.SetupLogger(cfg.Env)
	if err != nil {
		return nil, fmt.Errorf("error setting up logger: %v", err)
	}

	service, serviceDeps, err := newServerFunc(ctx, logger)
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
			interceptor.LoggerUnaryServerInterceptor(logger),
			verification.Interceptor(logger, verifier),
			interceptor.ValidationInterceptor(logger),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	registerServerFunc(grpcServer, service)

	gwMux := runtime.NewServeMux()

	err = registerEndPointFunc(
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
		Handler: otelhttp.NewHandler(corsHandler, "http-gateway"),
	}

	server := &grpcAppWrapper{
		app:          newAppFunc(serviceDeps),
		logger:       logger,
		serverConfig: cfg.Server,
		httpServer:   httpServer,
		grpcServer:   grpcServer,
		verifier:     verifier,
	}
	return server, nil
}

func (w *grpcAppWrapper) start() error {
	go func() {
		w.app.start()
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

func (w *grpcAppWrapper) shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, w.serverConfig.ShutdownTimeout)
	defer cancel()

	w.grpcServer.GracefulStop()
	w.verifier.Close()

	if err := w.httpServer.Shutdown(shutdownCtx); err != nil {
		w.logger.Error("HTTP server shutdown error", sl.Err(err))
	}

	return w.app.shutdown(shutdownCtx)
}

func selectVerifier(ssoConfig config.SSO) (verification.Verifier, error) {
	if ssoConfig.UseLocal {
		return verification.NewStubVerifier(ssoConfig.LocalUserId), nil
	}
	return verification.NewGRPCVerifier(ssoConfig.HostAddress)
}

func connectToDbPool(ctx context.Context, database config.Database) (*pgxpool.Pool, error) {
	dns := os.Getenv(database.ConnectionStringKey)
	dbPool, err := postgres.NewPostgresPool(ctx, dns, database.PingDuration)
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
