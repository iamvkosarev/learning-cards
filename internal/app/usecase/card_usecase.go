package usecase

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/domain/contracts"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"google.golang.org/grpc/codes"
)

type CardsUseCaseDeps struct {
	AuthVerifier contracts.AuthVerifier
	CardReader   contracts.CardReader
	CardWriter   contracts.CardWriter
	GroupReader  contracts.GroupReader
}

type CardsUseCase struct {
	deps CardsUseCaseDeps
}

func NewCardsUseCase(deps CardsUseCaseDeps) *CardsUseCase {
	return &CardsUseCase{
		deps: deps,
	}
}

func (c *CardsUseCase) Get(ctx context.Context, cardId entity.CardId) (entity.Card, error) {
	op := "usecase.CardsUseCase.Get"

	userId, err := c.deps.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return entity.Card{}, err
	}

	card, err := c.deps.CardReader.Get(ctx, cardId)
	if err != nil {
		return entity.Card{}, err
	}

	group, err := c.deps.GroupReader.Get(ctx, card.GroupId)
	if err != nil {
		return entity.Card{}, err
	}

	if err := checkViewGroupAccess(userId, group, op); err != nil {
		return entity.Card{}, err
	}

	return card, nil

}

func (c *CardsUseCase) List(ctx context.Context, groupId entity.GroupId) ([]entity.Card, error) {
	op := "usecase.CardsUseCase.List"

	userId, err := c.deps.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}
	group, err := c.deps.GroupReader.Get(ctx, groupId)
	if err != nil {
		return nil, err
	}

	if err := checkViewGroupAccess(userId, group, op); err != nil {
		return nil, err
	}
	cards, err := c.deps.CardReader.List(ctx, groupId)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (c *CardsUseCase) Create(ctx context.Context, groupId entity.GroupId, frontText, backText string) (
	entity.CardId,
	error,
) {
	op := "usecase.CardsUseCase.Create"

	userId, err := c.deps.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return 0, err
	}
	group, err := c.deps.GroupReader.Get(ctx, groupId)
	if err != nil {
		return 0, err
	}

	if userId != group.OwnerId {
		message := fmt.Sprintf("%v: user (id:%v) not owner of card groups", op, userId)
		return 0, entity.NewVerificationError(message, codes.PermissionDenied)
	}

	card := entity.Card{
		GroupId:   groupId,
		FrontText: frontText,
		BackText:  backText,
	}
	cardId, err := c.deps.CardWriter.Add(ctx, card)
	if err != nil {
		return 0, err
	}
	return cardId, nil
}

func (c *CardsUseCase) Update(ctx context.Context, updateCard entity.UpdateCard) error {
	op := "usecase.GroupUseCase.Update"

	userId, err := c.deps.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return err
	}
	card, err := c.deps.CardReader.Get(ctx, updateCard.Id)
	if err != nil {
		return err
	}

	group, err := c.deps.GroupReader.Get(ctx, card.GroupId)

	if err != nil {
		return err
	}

	if err := checkEditGroupAccess(userId, group, op); err != nil {
		return err
	}

	card.FrontText = updateCard.FrontText
	card.BackText = updateCard.BackText

	err = c.deps.CardWriter.Update(ctx, card)
	if err != nil {
		return err
	}

	return nil
}

func (c *CardsUseCase) Delete(ctx context.Context, id entity.CardId) error {
	op := "usecase.GroupUseCase.Delete"

	userId, err := c.deps.AuthVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return err
	}
	card, err := c.deps.CardReader.Get(ctx, id)
	if err != nil {
		return err
	}

	group, err := c.deps.GroupReader.Get(ctx, card.GroupId)

	if err != nil {
		return err
	}

	if err := checkEditGroupAccess(userId, group, op); err != nil {
		return err
	}

	err = c.deps.CardWriter.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
