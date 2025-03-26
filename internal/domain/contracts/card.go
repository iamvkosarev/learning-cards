package contracts

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type CardReader interface {
	Get(ctx context.Context, cardId entity.CardId) (entity.Card, error)
	List(ctx context.Context, groupId entity.GroupId) ([]entity.Card, error)
}

type CardWriter interface {
	Add(ctx context.Context, card entity.Card) (entity.CardId, error)
	Update(ctx context.Context, card entity.Card) error
	Delete(ctx context.Context, cardId entity.CardId) error
}
