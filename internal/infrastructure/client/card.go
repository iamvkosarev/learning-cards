package client

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/server/interceptor/verification"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

func (c *CardsClient) GetCard(ctx context.Context, cardId entity.CardId) (entity.Card, error) {
	resp, err := c.cardService.GetCard(ctx, &pb.GetCardRequest{CardId: int64(cardId)})
	if err != nil {
		return entity.Card{}, handleServiceErr(err)
	}
	createAt, err := time.Parse(time.RFC3339, resp.Card.CreatedAt)
	if err != nil {
		return entity.Card{}, fmt.Errorf("error parsing create time %w", err)
	}
	return entity.Card{
		Id:         entity.CardId(resp.Card.Id),
		GroupId:    entity.GroupId(resp.Card.GroupId),
		FrontText:  resp.Card.FrontText,
		BackText:   resp.Card.BackText,
		CreateTime: createAt,
	}, nil
}

func (c *CardsClient) ListCards(ctx context.Context, groupId entity.GroupId) ([]entity.Card, error) {
	resp, err := c.cardService.ListCards(ctx, &pb.ListCardsRequest{GroupId: int64(groupId)})
	if err != nil {
		return []entity.Card{}, handleServiceErr(err)
	}
	cards := make([]entity.Card, len(resp.Cards))
	for i, card := range resp.Cards {
		createAt, err := time.Parse(time.RFC3339, card.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error parsing create time %w", err)
		}
		cards[i] = entity.Card{
			Id:         entity.CardId(card.Id),
			GroupId:    entity.GroupId(card.GroupId),
			FrontText:  card.FrontText,
			BackText:   card.BackText,
			CreateTime: createAt,
		}
	}
	return cards, nil
}

func (c *CardsClient) GetGroup(ctx context.Context, groupId entity.GroupId) (entity.Group, error) {
	resp, err := c.cardService.GetGroup(ctx, &pb.GetGroupRequest{GroupId: int64(groupId)})
	if err != nil {
		return entity.Group{}, handleServiceErr(err)
	}
	if resp.Group == nil {
		return entity.Group{}, fmt.Errorf("error getting group from card service: group (groupId: %v) is nil", groupId)
	}
	createAt, err := time.Parse(time.RFC3339, resp.Group.CreatedAt)
	if err != nil {
		return entity.Group{}, fmt.Errorf("error parsing create time %w", err)
	}
	return entity.Group{
		Id:          entity.GroupId(resp.Group.Id),
		OwnerId:     entity.UserId(resp.Group.OwnerId),
		Name:        resp.Group.Name,
		Description: resp.Group.Description,
		CreateTime:  createAt,
		Visibility:  entity.GroupVisibility(resp.Group.Visibility),
	}, nil
}

func (c *CardsClient) ListGroups(ctx context.Context, _ entity.UserId) ([]entity.Group, error) {
	resp, err := c.cardService.ListGroups(ctx, &pb.ListGroupsRequest{})
	if err != nil {
		return []entity.Group{}, handleServiceErr(err)
	}
	groups := make([]entity.Group, len(resp.Groups))
	for i, group := range resp.Groups {
		createAt, err := time.Parse(time.RFC3339, group.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error parsing create time %w", err)
		}
		groups[i] = entity.Group{
			Id:          entity.GroupId(group.Id),
			OwnerId:     entity.UserId(group.OwnerId),
			Name:        group.Name,
			Description: group.Description,
			CreateTime:  createAt,
			Visibility:  entity.GroupVisibility(group.Visibility),
		}
	}
	return groups, nil
}

func handleServiceErr(err error) error {
	if statusErr, ok := status.FromError(err); ok {
		switch statusErr.Code() {
		case codes.PermissionDenied:
			return entity.NewVerificationError(statusErr.Err())
		case codes.NotFound:
			return entity.ErrGroupNotFound
		}
		return statusErr.Err()
	}
	return err
}
