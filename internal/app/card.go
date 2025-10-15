package app

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/server"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/server/interceptor/verification"
	"github.com/iamvkosarev/learning-cards/internal/model"
	"github.com/iamvkosarev/learning-cards/internal/module"
	"github.com/iamvkosarev/learning-cards/internal/otel/tracing"
	"github.com/iamvkosarev/learning-cards/internal/repository"
	"github.com/iamvkosarev/learning-cards/internal/service/japanese"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/sdk/trace"
	"log/slog"
)

type cardsAppDeps struct {
	dbPool        *pgxpool.Pool
	traceProvider *trace.TracerProvider
}
type cardsApp struct {
	*cardsAppDeps
}

func NewCardsApp(ctx context.Context, cfg *config.CardsConfig) (*grpcAppWrapper, error) {
	return buildApp(
		ctx, cfg.Common,
		func(ctx context.Context, logger *slog.Logger) (pb.CardServiceServer, *cardsAppDeps, error) {
			return prepareCardServer(ctx, cfg, logger)
		}, pb.RegisterCardServiceServer, pb.RegisterCardServiceHandlerFromEndpoint, newCardsApp,
	)
}

func newCardsApp(
	serviceDeps *cardsAppDeps,
) *cardsApp {
	return &cardsApp{
		cardsAppDeps: serviceDeps,
	}
}

func (c *cardsApp) start() error {
	return nil
}

func (c *cardsApp) shutdown(ctx context.Context) error {
	c.cardsAppDeps.dbPool.Close()
	return c.traceProvider.Shutdown(ctx)
}

func prepareCardServer(ctx context.Context, cfg *config.CardsConfig, logger *slog.Logger) (
	pb.CardServiceServer,
	*cardsAppDeps, error,
) {
	traceProvider, err := tracing.SetupTracingProvider(ctx, cfg.Tracing)
	if err != nil {
		return nil, nil, err
	}

	dbPool, err := connectToDbPool(ctx, cfg.Database)
	if err != nil {
		return nil, nil, err
	}

	rdb := redis.NewClient(
		&redis.Options{
			Addr: cfg.Redis.Endpoint,
		},
	)

	groupRepo := repository.NewGroupRepository(dbPool)
	userRepo := repository.NewUserRepository(dbPool)
	cardRepo := repository.NewCardRepository(dbPool)

	japaneseReader := japanese.NewReader(cfg.JapaneseReading, logger)

	verifier := server.VerifyFunc(verification.GetUserId)

	groupsService := module.NewGroups(
		module.GroupsDeps{
			GroupReader:  groupRepo,
			GroupWriter:  groupRepo,
			UserReader:   userRepo,
			UserWriter:   userRepo,
			UserVerifier: verifier,
		},
	)

	cardDecorator := module.NewDecorator(
		module.DecoratorDeps{
			GroupReader: groupRepo,
			CardReadingProviders: map[model.CardSideType]module.CardReadingProvider{
				model.CARD_SIDE_TYPE_JAPANESE: japaneseReader,
			},
		},
	)

	cardsService := module.NewCards(
		module.CardsDeps{
			CardWriter:         cardRepo,
			CardReader:         cardRepo,
			GroupAccessChecker: groupsService,
			GroupReader:        groupRepo,
			CardDecorator:      cardDecorator,
			Rdb:                rdb,
		},
	)

	server := server.NewCardServer(
		server.CardServerDeps{
			GroupService: groupsService,
			CardsService: cardsService,
			AuthVerifier: verifier,
			Logger:       logger,
		},
	)

	appDeps := &cardsAppDeps{
		dbPool:        dbPool,
		traceProvider: traceProvider,
	}
	return server, appDeps, nil
}
