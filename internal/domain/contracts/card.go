package contracts

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type CardReader interface {
	GetCard(ctx context.Context, cardId entity.CardId) (entity.Card, error)
	ListCards(ctx context.Context, groupId entity.GroupId) ([]entity.Card, error)
}

type CardWriter interface {
	AddCard(ctx context.Context, card entity.Card) (entity.CardId, error)
	UpdateCard(ctx context.Context, card entity.Card) error
	DeleteCard(ctx context.Context, cardId entity.CardId) error
}
