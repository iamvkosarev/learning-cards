package module

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const cardsTraceName = "module.cards"

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
	CheckReadGroupAccess(ctx context.Context, groupId model.GroupId) (*model.Group, error)
	CheckWriteGroupAccess(ctx context.Context, groupId model.GroupId) (*model.Group, error)
}

//go:generate minimock -i CardDecorator -o ./mocks/card_decorator_mock.go -n CardDecoratorMock -p mocks
type CardDecorator interface {
	DecorateCard(ctx context.Context, card *model.Card, group *model.Group) error
}

type CardsDeps struct {
	CardReader         CardReader
	CardWriter         CardWriter
	GroupReader        GroupReader
	GroupAccessChecker GroupAccessChecker
	CardDecorator      CardDecorator
}

type Cards struct {
	CardsDeps
	tracer trace.Tracer
}

func NewCards(deps CardsDeps) *Cards {
	return &Cards{
		CardsDeps: deps,
		tracer:    otel.Tracer(cardsTraceName),
	}
}

func (c *Cards) AddCard(
	ctx context.Context, groupId model.GroupId, sidesText []string,
) (
	model.CardId,
	error,
) {
	ctx, span := c.tracer.Start(ctx, "AddCard")
	defer span.End()
	if _, err := c.GroupAccessChecker.CheckWriteGroupAccess(ctx, groupId); err != nil {
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
	ctx, span := c.tracer.Start(ctx, "GetCard")
	defer span.End()
	card, err := c.CardReader.GetCard(ctx, cardId)
	if err != nil {
		return nil, err
	}

	group, err := c.GroupAccessChecker.CheckReadGroupAccess(ctx, card.GroupId)
	if err != nil {
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
	ctx, span := c.tracer.Start(ctx, "ListCards")
	defer span.End()
	group, err := c.GroupAccessChecker.CheckReadGroupAccess(ctx, groupId)
	if err != nil {
		return nil, err
	}

	cards, err := c.CardReader.ListCards(ctx, groupId)
	if err != nil {
		return nil, err
	}

	decorateCtx, span := c.tracer.Start(ctx, "DecorateCards")
	defer span.End()
	for _, card := range cards {
		if err = c.CardDecorator.DecorateCard(ctx, card); err != nil {
			return nil, err
		}
	}

	return cards, nil
}

func (c *Cards) UpdateCard(ctx context.Context, updateCard model.UpdateCard) error {
	ctx, span := c.tracer.Start(ctx, "UpdateCard")
	defer span.End()
	card, err := c.CardReader.GetCard(ctx, updateCard.Id)
	if err != nil {
		return err
	}

	if _, err = c.GroupAccessChecker.CheckWriteGroupAccess(ctx, card.GroupId); err != nil {
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
	ctx, span := c.tracer.Start(ctx, "DeleteCard")
	defer span.End()
	card, err := c.CardReader.GetCard(ctx, id)
	if err != nil {
		return err
	}

	if _, err = c.GroupAccessChecker.CheckWriteGroupAccess(ctx, card.GroupId); err != nil {
		return err
	}

	err = c.CardWriter.DeleteCard(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
