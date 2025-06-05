package service

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

//go:generate minimock -i CardReader -o ./mocks/card_reader_mock.go -n CardReaderMock -p mocks
type CardReader interface {
	GetCard(ctx context.Context, cardId entity.CardId) (entity.Card, error)
	ListCards(ctx context.Context, groupId entity.GroupId) ([]entity.Card, error)
}

//go:generate minimock -i CardWriter -o ./mocks/card_writer_mock.go -n CardWriterMock -p mocks
type CardWriter interface {
	AddCard(ctx context.Context, card entity.Card) (entity.CardId, error)
	UpdateCard(ctx context.Context, card entity.Card) error
	DeleteCard(ctx context.Context, cardId entity.CardId) error
}

//go:generate minimock -i GroupAccessChecker -o ./mocks/group_access_checker_mock.go -n GroupAccessChecker -p mocks
type GroupAccessChecker interface {
	CheckReadGroupAccess(ctx context.Context, groupId entity.GroupId) error
	CheckWriteGroupAccess(ctx context.Context, groupId entity.GroupId) error
}

type CardsServiceDeps struct {
	CardReader         CardReader
	CardWriter         CardWriter
	GroupAccessChecker GroupAccessChecker
}

type CardsService struct {
	CardsServiceDeps
}

func NewCardsService(deps CardsServiceDeps) *CardsService {
	return &CardsService{
		CardsServiceDeps: deps,
	}
}

func (c *CardsService) AddCard(
	ctx context.Context, groupId entity.GroupId, frontText,
	backText string,
) (
	entity.CardId,
	error,
) {
	if err := c.GroupAccessChecker.CheckWriteGroupAccess(ctx, groupId); err != nil {
		return 0, err
	}

	card := entity.Card{
		GroupId:   groupId,
		FrontText: frontText,
		BackText:  backText,
	}
	cardId, err := c.CardWriter.AddCard(ctx, card)
	if err != nil {
		return 0, err
	}
	return cardId, nil
}

func (c *CardsService) GetCard(ctx context.Context, cardId entity.CardId) (entity.Card, error) {
	card, err := c.CardReader.GetCard(ctx, cardId)
	if err != nil {
		return entity.Card{}, err
	}

	if err = c.GroupAccessChecker.CheckReadGroupAccess(ctx, card.GroupId); err != nil {
		return entity.Card{}, err
	}

	return card, nil

}

func (c *CardsService) ListCards(ctx context.Context, groupId entity.GroupId) (
	[]entity.Card,
	error,
) {
	if err := c.GroupAccessChecker.CheckReadGroupAccess(ctx, groupId); err != nil {
		return nil, err
	}

	cards, err := c.CardReader.ListCards(ctx, groupId)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (c *CardsService) UpdateCard(ctx context.Context, updateCard entity.UpdateCard) error {
	card, err := c.CardReader.GetCard(ctx, updateCard.Id)
	if err != nil {
		return err
	}

	if err = c.GroupAccessChecker.CheckWriteGroupAccess(ctx, card.GroupId); err != nil {
		return err
	}

	card.FrontText = updateCard.FrontText
	card.BackText = updateCard.BackText

	err = c.CardWriter.UpdateCard(ctx, card)
	if err != nil {
		return err
	}

	return nil
}

func (c *CardsService) DeleteCard(ctx context.Context, id entity.CardId) error {
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
