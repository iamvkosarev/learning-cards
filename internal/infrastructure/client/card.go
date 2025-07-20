package client

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/server/interceptor/verification"
	"github.com/iamvkosarev/learning-cards/internal/model"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

const healthCheckKey = "health-check"

type CardsClient struct {
	cardService pb.CardServiceClient
	conn        *grpc.ClientConn
}

func NewCardsClient(ctx context.Context, address string) (*CardsClient, error) {
	cardsConn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(AuthInterceptor()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating client to cards service: %w", err)
	}
	cardsService := pb.NewCardServiceClient(cardsConn)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ctx = context.WithValue(
		ctx, healthCheckKey, struct{}{},
	)

	_, err = cardsService.HealthCheck(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("error health check cards service: %w", err)
	}

	return &CardsClient{conn: cardsConn, cardService: cardsService}, nil
}

func AuthInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		op := "client.cards.AuthInterceptor"
		if value := ctx.Value(healthCheckKey); value != nil {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		token, err := verification.GetTokenFormContext(ctx)
		if err != nil {
			return fmt.Errorf("%s: could not get token: %w", op, err)
		}

		if token != "" {
			md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", token))
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
func (c *CardsClient) Close() {
	c.conn.Close()
}

func (c *CardsClient) GetCard(ctx context.Context, cardId model.CardId) (*model.Card, error) {
	resp, err := c.cardService.GetCard(ctx, &pb.GetCardRequest{CardId: int64(cardId)})
	if err != nil {
		return nil, err
	}
	createAt, err := time.Parse(time.RFC3339, resp.Card.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error parsing create time %w", err)
	}

	return &model.Card{
		Id:         model.CardId(resp.Card.Id),
		GroupId:    model.GroupId(resp.Card.GroupId),
		Sides:      parseCardSides(resp.Card.Sides),
		CreateTime: createAt,
	}, nil
}

func (c *CardsClient) ListCards(ctx context.Context, groupId model.GroupId) ([]*model.Card, error) {
	resp, err := c.cardService.ListCards(ctx, &pb.ListCardsRequest{GroupId: int64(groupId)})
	if err != nil {
		return nil, err
	}
	cards := make([]*model.Card, len(resp.Cards))
	for i, card := range resp.Cards {
		var createAt time.Time
		createAt, err = time.Parse(time.RFC3339, card.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error parsing create time %w", err)
		}
		cards[i] = &model.Card{
			Id:         model.CardId(card.Id),
			GroupId:    model.GroupId(card.GroupId),
			Sides:      parseCardSides(card.Sides),
			CreateTime: createAt,
		}
	}
	return cards, nil
}

func parseCardSides(sides []*pb.CardSide) []model.CardSide {
	parsedSides := make([]model.CardSide, len(sides))
	for i, side := range sides {
		readingPairs := make([]model.ReadingPair, len(side.ReadingPairs))
		for j, readingPair := range side.ReadingPairs {
			readingPairs[j] = model.ReadingPair{
				Text:    readingPair.Text,
				Reading: readingPair.Reading,
			}
		}
		parsedSides[i] = model.CardSide{
			Text:         side.Text,
			ReadingPairs: readingPairs,
		}
	}
	return parsedSides
}

func (c *CardsClient) GetGroup(ctx context.Context, groupId model.GroupId) (*model.Group, error) {
	resp, err := c.cardService.GetGroup(ctx, &pb.GetGroupRequest{GroupId: int64(groupId)})
	if err != nil {
		return nil, err
	}
	if resp.Group == nil {
		return nil, fmt.Errorf("error getting group from card service: group (groupId: %v) is nil", groupId)
	}
	createAt, err := time.Parse(time.RFC3339, resp.Group.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error parsing create time %w", err)
	}
	return &model.Group{
		Id:          model.GroupId(resp.Group.Id),
		OwnerId:     model.UserId(resp.Group.OwnerId),
		Name:        resp.Group.Name,
		Description: resp.Group.Description,
		CreateTime:  createAt,
		Visibility:  model.GroupVisibility(resp.Group.Visibility),
		CardSideTypes: []model.CardSideType{
			model.CardSideType(resp.Group.CardSideTypes[0]),
			model.CardSideType(resp.Group.CardSideTypes[1]),
		},
	}, nil
}

func (c *CardsClient) ListGroups(ctx context.Context, _ model.UserId) ([]*model.Group, error) {
	resp, err := c.cardService.ListGroups(ctx, &pb.ListGroupsRequest{})
	if err != nil {
		return nil, err
	}
	groups := make([]*model.Group, len(resp.Groups))
	for i, group := range resp.Groups {
		createAt, err := time.Parse(time.RFC3339, group.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error parsing create time %w", err)
		}
		groups[i] = &model.Group{
			Id:          model.GroupId(group.Id),
			OwnerId:     model.UserId(group.OwnerId),
			Name:        group.Name,
			Description: group.Description,
			CreateTime:  createAt,
			Visibility:  model.GroupVisibility(group.Visibility),
		}
	}
	return groups, nil
}
