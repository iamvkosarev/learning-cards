package module

import (
	"context"
	"github.com/iamvkosarev/go-shared-utils/logger/sl"
	"github.com/iamvkosarev/learning-cards/internal/model"
	"log/slog"
)

//go:generate minimock -i CardReadingProvider -o ./mocks/card_reading_provider_mock.go -n CardReadingProviderMock -p mocks
type CardReadingProvider interface {
	GetCardReading(ctx context.Context, text string) ([]model.ReadingPair, error)
}
type DecoratorDeps struct {
	GroupReader          GroupReader
	CardReadingProviders map[model.CardSideType]CardReadingProvider
	Logger               *slog.Logger
}

type Decorator struct {
	DecoratorDeps
}

func NewDecorator(deps DecoratorDeps) *Decorator {
	return &Decorator{
		DecoratorDeps: deps,
	}
}

func (d *Decorator) TryDecorateCard(ctx context.Context, card *model.Card, group *model.Group) {
	var err error
	for i, side := range card.Sides {
		switch group.CardSideTypes[i] {
		case model.CARD_SIDE_TYPE_JAPANESE:
			card.Sides[i], err = d.decorateJapaneseCard(ctx, side)
			if err != nil {
				d.Logger.Error(
					"failed to decorate card", slog.String("side-type", "japanese"), slog.String(
						"text", side.Text,
					), sl.Err(err),
				)
			}
		}
	}
}

func (d *Decorator) decorateJapaneseCard(ctx context.Context, side model.CardSide) (
	model.CardSide,
	error,
) {
	if readingProvider, ok := d.CardReadingProviders[model.CARD_SIDE_TYPE_JAPANESE]; ok {
		readingPairs, err := readingProvider.GetCardReading(ctx, side.Text)
		if err != nil {
			return model.CardSide{}, err
		}
		side.ReadingPairs = readingPairs
	}
	return side, nil
}
