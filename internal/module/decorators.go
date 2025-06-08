package module

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/model"
)

//go:generate minimock -i CardReadingProvider -o ./mocks/card_reading_provider_mock.go -n CardReadingProviderMock -p mocks
type CardReadingProvider interface {
	GetCardReading(ctx context.Context, text string) ([]model.ReadingPair, error)
}
type DecoratorDeps struct {
	GroupReader          GroupReader
	CardReadingProviders map[model.CardSideType]CardReadingProvider
}

type Decorator struct {
	DecoratorDeps
}

func NewDecorator(deps DecoratorDeps) *Decorator {
	return &Decorator{
		DecoratorDeps: deps,
	}
}

func (c *Decorator) DecorateCard(ctx context.Context, card *model.Card) error {
	group, err := c.GroupReader.GetGroup(ctx, card.GroupId)
	if err != nil {
		return err
	}
	for i, side := range card.Sides {
		switch group.CardSideTypes[i] {
		case model.CARD_SIDE_TYPE_JAPANESE:
			card.Sides[i], err = c.decorateJapaneseCard(ctx, side)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Decorator) decorateJapaneseCard(ctx context.Context, side model.CardSide) (
	model.CardSide,
	error,
) {
	if readingProvider, ok := c.CardReadingProviders[model.CARD_SIDE_TYPE_JAPANESE]; ok {
		readingPairs, err := readingProvider.GetCardReading(ctx, side.Text)
		if err != nil {
			return model.CardSide{}, err
		}
		side.ReadingPairs = readingPairs
	}
	return side, nil
}
