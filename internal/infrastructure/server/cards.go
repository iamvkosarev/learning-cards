package server

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/domain/contracts"
	server "github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/grpc/interceptor/verification"
	sqlRepository "github.com/iamvkosarev/learning-cards/internal/infrastructure/repository/postgres"
	usecase "github.com/iamvkosarev/learning-cards/internal/usecase"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type CardsServiceDeps struct {
	dbPool *pgxpool.Pool
}
type CardsServer struct {
	CardsServiceDeps
}

func BuildCardsServer(ctx context.Context, cfg *config.CardsConfig) (*Wrapper, error) {
	return BuildServer(
		ctx, cfg.Common,
		func(ctx context.Context, logger *slog.Logger) (pb.CardServiceServer, CardsServiceDeps, error) {
			return prepareCardService(ctx, cfg, logger)
		}, pb.RegisterCardServiceServer, pb.RegisterCardServiceHandlerFromEndpoint, NewCardsServer,
	)
}

func NewCardsServer(
	serviceDeps CardsServiceDeps,
) *CardsServer {
	return &CardsServer{
		CardsServiceDeps: serviceDeps,
	}
}

func (c *CardsServer) Start() error {
	return nil
}

func (c *CardsServer) Shutdown(context.Context) {
	c.CardsServiceDeps.dbPool.Close()
}

func prepareCardService(ctx context.Context, cfg *config.CardsConfig, logger *slog.Logger) (
	pb.CardServiceServer,
	CardsServiceDeps, error,
) {
	dbPool, err := ConnectToDbPool(ctx, cfg.Database)
	if err != nil {
		return nil, CardsServiceDeps{}, err
	}

	groupRepo := sqlRepository.NewGroupRepository(dbPool)
	userRepo := sqlRepository.NewUserRepository(dbPool)
	cardRepo := sqlRepository.NewCardRepository(dbPool)

	verifier := contracts.VerifyFunc(verification.GetUserId)

	var groupUseCase server.GroupUseCase = usecase.NewGroupUseCase(
		usecase.GroupUseCaseDeps{
			GroupReader: groupRepo,
			GroupWriter: groupRepo,
			UserReader:  userRepo,
			UserWriter:  userRepo,
		},
	)
	var cardsUseCase server.CardsUseCase = usecase.NewCardsUseCase(
		usecase.CardsUseCaseDeps{
			GroupReader: groupRepo,
			CardWriter:  cardRepo,
			CardReader:  cardRepo,
		},
	)

	cardsService := server.NewCardService(
		server.CardServiceDeps{
			GroupUseCase: groupUseCase,
			CardsUseCase: cardsUseCase,
			AuthVerifier: verifier,
			Logger:       logger,
		},
	)

	serviceDeps := CardsServiceDeps{
		dbPool: dbPool,
	}
	return cardsService, serviceDeps, nil
}
