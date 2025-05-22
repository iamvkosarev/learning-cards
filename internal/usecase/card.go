package usecase

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/domain/contracts"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CardsUseCaseDeps struct {
	CardReader   contracts.CardReader
	CardWriter   contracts.CardWriter
	GroupReader  contracts.GroupReader
	AuthVerifier contracts.AuthVerifier
}

type CardsUseCase struct {
	CardsUseCaseDeps
}

func NewCardsUseCase(deps CardsUseCaseDeps) *CardsUseCase {
	return &CardsUseCase{
		CardsUseCaseDeps: deps,
	}
}

func (c *CardsUseCase) GetCard(ctx context.Context, cardId entity.CardId) (entity.Card, error) {
	op := "usecase.CardsUseCase.GetCard"

	userId, err := c.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return entity.Card{}, err
	}

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

func (c *CardsUseCase) ListCards(ctx context.Context, groupId entity.GroupId) ([]entity.Card, error) {
	op := "usecase.CardsUseCase.ListCards"

	userId, err := c.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
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

func (c *CardsUseCase) Create(ctx context.Context, groupId entity.GroupId, frontText, backText string) (
	entity.CardId,
	error,
) {
	op := "usecase.CardsUseCase.CreateGroup"

	userId, err := c.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return 0, err
	}

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

func (c *CardsUseCase) UpdateCard(ctx context.Context, updateCard entity.UpdateCard) error {
	op := "usecase.GroupUseCase.UpdateCard"

	userId, err := c.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return err
	}

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

func (c *CardsUseCase) DeleteCard(ctx context.Context, id entity.CardId) error {
	op := "usecase.GroupUseCase.DeleteCard"

	userId, err := c.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return err
	}

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
