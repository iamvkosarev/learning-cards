package client

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	pb "github.com/iamvkosarev/learning-cards/pkg/proto/learning_cards/v1"
)

type CardsClient struct {
	pb.CardServiceClient
}

func NewCardsClient() *CardsClient {
	return &CardsClient{}
}

func (c CardsClient) GetCard(ctx context.Context, cardId entity.CardId) (entity.Card, error) {
	//TODO implement me
	panic("implement me")
}

func (c CardsClient) ListCards(ctx context.Context, groupId entity.GroupId) ([]entity.Card, error) {
	//TODO implement me
	panic("implement me")
}

func (c CardsClient) GetGroup(ctx context.Context, groupId entity.GroupId) (entity.Group, error) {
	//TODO implement me
	panic("implement me")
}

func (c CardsClient) ListGroups(ctx context.Context, id entity.UserId) ([]entity.Group, error) {
	//TODO implement me
	panic("implement me")
}
