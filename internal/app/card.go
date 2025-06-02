package app

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/config"
	sqlRepository "github.com/iamvkosarev/learning-cards/internal/infrastructure/repository/postgres"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/server"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/server/interceptor/verification"
	"github.com/iamvkosarev/learning-cards/internal/service"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type cardsAppDeps struct {
	dbPool *pgxpool.Pool
}
type cardsApp struct {
	cardsAppDeps
}

func NewCardsApp(ctx context.Context, cfg *config.CardsConfig) (*grpcAppWrapper, error) {
	return buildApp(
		ctx, cfg.Common,
		func(ctx context.Context, logger *slog.Logger) (pb.CardServiceServer, cardsAppDeps, error) {
			return prepareCardServer(ctx, cfg, logger)
		}, pb.RegisterCardServiceServer, pb.RegisterCardServiceHandlerFromEndpoint, newCardsApp,
	)
}

func newCardsApp(
	serviceDeps cardsAppDeps,
) *cardsApp {
	return &cardsApp{
		cardsAppDeps: serviceDeps,
	}
}

func (c *cardsApp) start() error {
	return nil
}

func (c *cardsApp) shutdown(context.Context) {
	c.cardsAppDeps.dbPool.Close()
}

func prepareCardServer(ctx context.Context, cfg *config.CardsConfig, logger *slog.Logger) (
	pb.CardServiceServer,
	cardsAppDeps, error,
) {
	dbPool, err := connectToDbPool(ctx, cfg.Database)
	if err != nil {
		return nil, cardsAppDeps{}, err
	}

	groupRepo := sqlRepository.NewGroupRepository(dbPool)
	userRepo := sqlRepository.NewUserRepository(dbPool)
	cardRepo := sqlRepository.NewCardRepository(dbPool)

	verifier := server.VerifyFunc(verification.GetUserId)

	var groupService server.GroupService = service.NewGroupService(
		service.GroupServiceDeps{
			GroupReader: groupRepo,
			GroupWriter: groupRepo,
			UserReader:  userRepo,
			UserWriter:  userRepo,
		},
	)
	var cardsService server.CardsService = service.NewCardsService(
		service.CardsServiceDeps{
			GroupReader: groupRepo,
			CardWriter:  cardRepo,
			CardReader:  cardRepo,
		},
	)

	server := server.NewCardServer(
		server.CardServerDeps{
			GroupService: groupService,
			CardsService: cardsService,
			AuthVerifier: verifier,
			Logger:       logger,
		},
	)

	appDeps := cardsAppDeps{
		dbPool: dbPool,
	}
	return server, appDeps, nil
}
