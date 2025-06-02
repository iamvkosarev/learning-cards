package service

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

type CardsServiceDeps struct {
	CardReader  CardReader
	CardWriter  CardWriter
	GroupReader GroupReader
}

type CardsService struct {
	CardsServiceDeps
}

func NewCardsService(deps CardsServiceDeps) *CardsService {
	return &CardsService{
		CardsServiceDeps: deps,
	}
}

func (c *CardsService) GetCard(ctx context.Context, userId entity.UserId, cardId entity.CardId) (entity.Card, error) {
	op := "service.CardsService.GetCard"

	card, err := c.CardReader.GetCard(ctx, cardId)
	if err != nil {
		return entity.Card{}, err
	}

	group, err := c.GroupReader.GetGroup(ctx, card.GroupId)
	if err != nil {
		return entity.Card{}, err
	}

	if err := checkViewGroupAccess(userId, group, op); err != nil {
		return entity.Card{}, err
	}

	return card, nil

}

func (c *CardsService) ListCards(ctx context.Context, userId entity.UserId, groupId entity.GroupId) (
	[]entity.Card,
	error,
) {
	op := "service.CardsService.ListCards"

	group, err := c.GroupReader.GetGroup(ctx, groupId)
	if err != nil {
		return nil, err
	}

	if err := checkViewGroupAccess(userId, group, op); err != nil {
		return nil, err
	}
	cards, err := c.CardReader.ListCards(ctx, groupId)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (c *CardsService) Create(
	ctx context.Context, userId entity.UserId, groupId entity.GroupId, frontText,
	backText string,
) (
	entity.CardId,
	error,
) {
	op := "service.CardsService.CreateGroup"

	group, err := c.GroupReader.GetGroup(ctx, groupId)
	if err != nil {
		return 0, err
	}

	if userId != group.OwnerId {
		message := fmt.Sprintf("%v: user (id:%v) not owner of card groups", op, userId)
		return 0, entity.NewVerificationError(status.Error(codes.PermissionDenied, message))
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

func (c *CardsService) UpdateCard(ctx context.Context, userId entity.UserId, updateCard entity.UpdateCard) error {
	op := "service.GroupService.UpdateCard"

	card, err := c.CardReader.GetCard(ctx, updateCard.Id)
	if err != nil {
		return err
	}

	group, err := c.GroupReader.GetGroup(ctx, card.GroupId)

	if err != nil {
		return err
	}

	if err := checkEditGroupAccess(userId, group, op); err != nil {
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

func (c *CardsService) DeleteCard(ctx context.Context, userId entity.UserId, id entity.CardId) error {
	op := "service.GroupService.DeleteCard"

	card, err := c.CardReader.GetCard(ctx, id)
	if err != nil {
		return err
	}

	group, err := c.GroupReader.GetGroup(ctx, card.GroupId)

	if err != nil {
		return err
	}

	if err := checkEditGroupAccess(userId, group, op); err != nil {
		return err
	}

	err = c.CardWriter.DeleteCard(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
