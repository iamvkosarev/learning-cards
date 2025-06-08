package module

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/model"
)

//go:generate minimock -i CardReader -o ./mocks/card_reader_mock.go -n CardReaderMock -p mocks
type CardReader interface {
	GetCard(ctx context.Context, cardId model.CardId) (*model.Card, error)
	ListCards(ctx context.Context, groupId model.GroupId) ([]*model.Card, error)
}

//go:generate minimock -i CardWriter -o ./mocks/card_writer_mock.go -n CardWriterMock -p mocks
type CardWriter interface {
	AddCard(ctx context.Context, card *model.Card) (model.CardId, error)
	UpdateCard(ctx context.Context, card *model.Card) error
	DeleteCard(ctx context.Context, cardId model.CardId) error
}

//go:generate minimock -i GroupAccessChecker -o ./mocks/group_access_checker_mock.go -n GroupAccessCheckerMock -p mocks
type GroupAccessChecker interface {
	CheckReadGroupAccess(ctx context.Context, groupId model.GroupId) error
	CheckWriteGroupAccess(ctx context.Context, groupId model.GroupId) error
}

//go:generate minimock -i CardDecorator -o ./mocks/card_decorator_mock.go -n CardDecoratorMock -p mocks
type CardDecorator interface {
	DecorateCard(ctx context.Context, card *model.Card) error
}

type CardsDeps struct {
	CardReader         CardReader
	CardWriter         CardWriter
	GroupAccessChecker GroupAccessChecker
	CardDecorator      CardDecorator
}

type Cards struct {
	CardsDeps
}

func NewCards(deps CardsDeps) *Cards {
	return &Cards{
		CardsDeps: deps,
	}
}

func (c *Cards) AddCard(
	ctx context.Context, groupId model.GroupId, sidesText []string,
) (
	model.CardId,
	error,
) {
	if err := c.GroupAccessChecker.CheckWriteGroupAccess(ctx, groupId); err != nil {
		return 0, err
	}

	card := &model.Card{
		GroupId: groupId,
		Sides:   make([]model.CardSide, len(sidesText)),
	}
	for i, text := range sidesText {
		card.Sides[i].Text = text
	}
	cardId, err := c.CardWriter.AddCard(ctx, card)
	if err != nil {
		return 0, err
	}

	return cardId, nil
}

func (c *Cards) GetCard(ctx context.Context, cardId model.CardId) (*model.Card, error) {
	card, err := c.CardReader.GetCard(ctx, cardId)
	if err != nil {
		return nil, err
	}

	if err = c.GroupAccessChecker.CheckReadGroupAccess(ctx, card.GroupId); err != nil {
		return nil, err
	}

	if err = c.CardDecorator.DecorateCard(ctx, card); err != nil {
		return nil, err
	}

	return card, nil

}

func (c *Cards) ListCards(ctx context.Context, groupId model.GroupId) (
	[]*model.Card,
	error,
) {
	if err := c.GroupAccessChecker.CheckReadGroupAccess(ctx, groupId); err != nil {
		return nil, err
	}

	cards, err := c.CardReader.ListCards(ctx, groupId)
	if err != nil {
		return nil, err
	}
	for _, card := range cards {
		if err = c.CardDecorator.DecorateCard(ctx, card); err != nil {
			return nil, err
		}
	}

	return cards, nil
}

func (c *Cards) UpdateCard(ctx context.Context, updateCard model.UpdateCard) error {
	card, err := c.CardReader.GetCard(ctx, updateCard.Id)
	if err != nil {
		return err
	}

	if err = c.GroupAccessChecker.CheckWriteGroupAccess(ctx, card.GroupId); err != nil {
		return err
	}

	for i, text := range updateCard.SidesText {
		if text == "" {
			continue
		}
		card.Sides[i].Text = text
	}

	err = c.CardWriter.UpdateCard(ctx, card)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cards) DeleteCard(ctx context.Context, id model.CardId) error {
	card, err := c.CardReader.GetCard(ctx, id)
	if err != nil {
		return err
	}

	if err = c.GroupAccessChecker.CheckWriteGroupAccess(ctx, card.GroupId); err != nil {
		return err
	}

	err = c.CardWriter.DeleteCard(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
