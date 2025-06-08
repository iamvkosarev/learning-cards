package module

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/iamvkosarev/learning-cards/internal/model"
	"github.com/iamvkosarev/learning-cards/internal/module/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestCardsService_AddCard(t *testing.T) {
	ctxNotCorrectUser := context.Background()
	md := metadata.Pairs("authorization", "Bearer correct-user-token")
	ctxCorrectUser := metadata.NewOutgoingContext(ctxNotCorrectUser, md)
	cardId := model.CardId(200)
	groupId := model.GroupId(200)
	tests := []struct {
		name      string
		ctx       context.Context
		groupId   model.GroupId
		frontText string
		backText  string
		result    model.CardId
		err       error
	}{
		{
			name:      "success",
			ctx:       ctxCorrectUser,
			groupId:   groupId,
			frontText: "Test Front Text",
			backText:  "Test Back Text",
			result:    cardId,
			err:       nil,
		},
		{
			name:    "no access",
			ctx:     ctxNotCorrectUser,
			groupId: groupId,
			err:     model.ErrGroupWriteAccessDenied,
		},
	}

	mc := minimock.NewController(t)

	cardDecoratorMock := mocks.NewCardDecoratorMock(mc)
	cardDecoratorMock.DecorateCardMock.Return(nil)

	cardsWriterMock := mocks.NewCardWriterMock(mc)
	cardsWriterMock.AddCardMock.Return(cardId, nil)

	groupAccessChecker := mocks.NewGroupAccessCheckerMock(mc)
	groupAccessChecker.CheckWriteGroupAccessMock.When(ctxCorrectUser, groupId).Then(nil)
	groupAccessChecker.CheckWriteGroupAccessMock.When(ctxNotCorrectUser, groupId).Then(model.ErrGroupWriteAccessDenied)
	cardService := NewCards(
		CardsDeps{
			CardWriter:         cardsWriterMock,
			GroupAccessChecker: groupAccessChecker,
			CardDecorator:      cardDecoratorMock,
		},
	)
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				newCardId, err := cardService.AddCard(
					test.ctx, test.groupId, test.frontText,
					test.backText,
				)
				if test.err == nil {
					assert.Equal(t, test.result, newCardId)
				}
				assert.Equal(t, test.err, err)
			},
		)
	}
}

func TestCardsService_DeleteCard(t *testing.T) {
	ctxNotCorrectUser := context.Background()
	md := metadata.Pairs("authorization", "Bearer correct-user-token")
	ctxCorrectUser := metadata.NewOutgoingContext(ctxNotCorrectUser, md)
	correctCardId := model.CardId(200)
	notExistsCardId := model.CardId(0)
	card := &model.Card{
		Id: correctCardId,
	}
	groupId := model.GroupId(0)
	tests := []struct {
		name   string
		ctx    context.Context
		cardId model.CardId
		err    error
	}{
		{
			name:   "success",
			ctx:    ctxCorrectUser,
			cardId: correctCardId,
			err:    nil,
		},
		{
			name:   "no access",
			ctx:    ctxNotCorrectUser,
			cardId: correctCardId,
			err:    model.ErrGroupWriteAccessDenied,
		},
		{
			name:   "no found",
			ctx:    ctxNotCorrectUser,
			cardId: notExistsCardId,
			err:    model.ErrCardNotFound,
		},
	}

	mc := minimock.NewController(t)

	cardsReaderMock := mocks.NewCardReaderMock(mc)
	cardsReaderMock.GetCardMock.When(minimock.AnyContext, correctCardId).Then(card, nil)
	cardsReaderMock.GetCardMock.When(minimock.AnyContext, notExistsCardId).Then(card, model.ErrCardNotFound)

	cardsWriterMock := mocks.NewCardWriterMock(mc)
	cardsWriterMock.DeleteCardMock.When(minimock.AnyContext, correctCardId).Then(nil)

	groupAccessChecker := mocks.NewGroupAccessCheckerMock(mc)
	groupAccessChecker.CheckWriteGroupAccessMock.When(ctxCorrectUser, groupId).Then(nil)
	groupAccessChecker.CheckWriteGroupAccessMock.When(ctxNotCorrectUser, groupId).Then(model.ErrGroupWriteAccessDenied)
	cardService := NewCards(
		CardsDeps{
			CardReader:         cardsReaderMock,
			CardWriter:         cardsWriterMock,
			GroupAccessChecker: groupAccessChecker,
		},
	)
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				err := cardService.DeleteCard(
					test.ctx, test.cardId,
				)
				assert.Equal(t, test.err, err)
			},
		)
	}
}

func TestCardsService_GetCard(t *testing.T) {
	ctxNotCorrectUser := context.Background()
	md := metadata.Pairs("authorization", "Bearer correct-user-token")
	ctxCorrectUser := metadata.NewOutgoingContext(ctxNotCorrectUser, md)
	correctCardId := model.CardId(200)
	notExistCardId := model.CardId(0)
	correctCard := &model.Card{
		Id: correctCardId,
	}
	groupId := model.GroupId(0)
	tests := []struct {
		name   string
		ctx    context.Context
		cardId model.CardId
		result *model.Card
		err    error
	}{
		{
			name:   "success",
			ctx:    ctxCorrectUser,
			cardId: correctCardId,
			result: correctCard,
			err:    nil,
		},
		{
			name:   "no access",
			ctx:    ctxNotCorrectUser,
			cardId: correctCardId,
			err:    model.ErrGroupWriteAccessDenied,
		},
		{
			name:   "no found",
			ctx:    ctxNotCorrectUser,
			cardId: notExistCardId,
			err:    model.ErrCardNotFound,
		},
	}

	mc := minimock.NewController(t)

	cardsReaderMock := mocks.NewCardReaderMock(mc)
	cardsReaderMock.GetCardMock.When(minimock.AnyContext, correctCardId).Then(correctCard, nil)
	cardsReaderMock.GetCardMock.When(minimock.AnyContext, notExistCardId).Then(correctCard, model.ErrCardNotFound)

	groupAccessChecker := mocks.NewGroupAccessCheckerMock(mc)
	groupAccessChecker.CheckReadGroupAccessMock.When(ctxCorrectUser, groupId).Then(nil)
	groupAccessChecker.CheckReadGroupAccessMock.When(ctxNotCorrectUser, groupId).Then(model.ErrGroupWriteAccessDenied)
	cardService := NewCards(
		CardsDeps{
			CardReader:         cardsReaderMock,
			GroupAccessChecker: groupAccessChecker,
		},
	)
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				card, err := cardService.GetCard(
					test.ctx, test.cardId,
				)
				if test.err == nil {
					assert.Equal(t, test.result, card)
				}
				assert.Equal(t, test.err, err)
			},
		)
	}
}
func TestCardsService_ListCards(t *testing.T) {
	ctxNotCorrectUser := context.Background()
	md := metadata.Pairs("authorization", "Bearer correct-user-token")
	ctx := metadata.NewOutgoingContext(ctxNotCorrectUser, md)
	groupId := model.GroupId(200)
	notExistGroupId := model.GroupId(1)
	correctCards := []*model.Card{
		{},
	}
	tests := []struct {
		name    string
		ctx     context.Context
		groupId model.GroupId
		result  []*model.Card
		err     error
	}{
		{
			name:    "success",
			ctx:     ctx,
			groupId: groupId,
			result:  correctCards,
			err:     nil,
		},
		{
			name:    "no access",
			ctx:     ctxNotCorrectUser,
			groupId: groupId,
			err:     model.ErrGroupWriteAccessDenied,
		},
		{
			name:    "no found",
			ctx:     ctx,
			groupId: notExistGroupId,
			err:     model.ErrGroupNotFound,
		},
	}

	mc := minimock.NewController(t)

	cardsReaderMock := mocks.NewCardReaderMock(mc)
	cardsReaderMock.ListCardsMock.When(minimock.AnyContext, groupId).Then(correctCards, nil)
	cardsReaderMock.ListCardsMock.When(minimock.AnyContext, notExistGroupId).Then(nil, model.ErrGroupNotFound)

	groupAccessChecker := mocks.NewGroupAccessCheckerMock(mc)
	groupAccessChecker.CheckReadGroupAccessMock.When(ctx, groupId).Then(nil)
	groupAccessChecker.CheckReadGroupAccessMock.When(ctx, notExistGroupId).Then(nil)
	groupAccessChecker.CheckReadGroupAccessMock.When(
		ctxNotCorrectUser, groupId,
	).Then(model.ErrGroupWriteAccessDenied)
	cardService := NewCards(
		CardsDeps{
			CardReader:         cardsReaderMock,
			GroupAccessChecker: groupAccessChecker,
		},
	)
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				cards, err := cardService.ListCards(
					test.ctx, test.groupId,
				)
				if test.err == nil {
					assert.Equal(t, test.result, cards)
				}
				assert.Equal(t, test.err, err)
			},
		)
	}
}
