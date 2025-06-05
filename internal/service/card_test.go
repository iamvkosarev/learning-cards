package service

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"github.com/iamvkosarev/learning-cards/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestCardsService_AddCard(t *testing.T) {
	ctxNotCorrectUser := context.Background()
	md := metadata.Pairs("authorization", "Bearer correct-user-token")
	ctxCorrectUser := metadata.NewOutgoingContext(ctxNotCorrectUser, md)
	cardId := entity.CardId(200)
	groupId := entity.GroupId(0)
	tests := []struct {
		name      string
		ctx       context.Context
		groupId   entity.GroupId
		frontText string
		backText  string
		result    entity.CardId
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
			err:     entity.ErrGroupWriteAccessDenied,
		},
	}

	mc := minimock.NewController(t)

	cardsWriterMock := mocks.NewCardWriterMock(mc)
	cardsWriterMock.AddCardMock.Return(cardId, nil)

	groupAccessChecker := mocks.NewGroupAccessChecker(mc)
	groupAccessChecker.CheckWriteGroupAccessMock.When(ctxCorrectUser, groupId).Then(nil)
	groupAccessChecker.CheckWriteGroupAccessMock.When(ctxNotCorrectUser, groupId).Then(entity.ErrGroupWriteAccessDenied)
	cardService := NewCardsService(
		CardsServiceDeps{
			CardWriter:         cardsWriterMock,
			GroupAccessChecker: groupAccessChecker,
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
	correctCardId := entity.CardId(200)
	notExistsCardId := entity.CardId(0)
	card := entity.Card{
		Id: correctCardId,
	}
	groupId := entity.GroupId(0)
	tests := []struct {
		name   string
		ctx    context.Context
		cardId entity.CardId
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
			err:    entity.ErrGroupWriteAccessDenied,
		},
		{
			name:   "no found",
			ctx:    ctxNotCorrectUser,
			cardId: notExistsCardId,
			err:    entity.ErrCardNotFound,
		},
	}

	mc := minimock.NewController(t)

	cardsReaderMock := mocks.NewCardReaderMock(mc)
	cardsReaderMock.GetCardMock.When(minimock.AnyContext, correctCardId).Then(card, nil)
	cardsReaderMock.GetCardMock.When(minimock.AnyContext, notExistsCardId).Then(card, entity.ErrCardNotFound)

	cardsWriterMock := mocks.NewCardWriterMock(mc)
	cardsWriterMock.DeleteCardMock.When(minimock.AnyContext, correctCardId).Then(nil)

	groupAccessChecker := mocks.NewGroupAccessChecker(mc)
	groupAccessChecker.CheckWriteGroupAccessMock.When(ctxCorrectUser, groupId).Then(nil)
	groupAccessChecker.CheckWriteGroupAccessMock.When(ctxNotCorrectUser, groupId).Then(entity.ErrGroupWriteAccessDenied)
	cardService := NewCardsService(
		CardsServiceDeps{
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
	correctCardId := entity.CardId(200)
	notExistCardId := entity.CardId(0)
	correctCard := entity.Card{
		Id: correctCardId,
	}
	groupId := entity.GroupId(0)
	tests := []struct {
		name   string
		ctx    context.Context
		cardId entity.CardId
		result entity.Card
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
			err:    entity.ErrGroupWriteAccessDenied,
		},
		{
			name:   "no found",
			ctx:    ctxNotCorrectUser,
			cardId: notExistCardId,
			err:    entity.ErrCardNotFound,
		},
	}

	mc := minimock.NewController(t)

	cardsReaderMock := mocks.NewCardReaderMock(mc)
	cardsReaderMock.GetCardMock.When(minimock.AnyContext, correctCardId).Then(correctCard, nil)
	cardsReaderMock.GetCardMock.When(minimock.AnyContext, notExistCardId).Then(correctCard, entity.ErrCardNotFound)

	groupAccessChecker := mocks.NewGroupAccessChecker(mc)
	groupAccessChecker.CheckReadGroupAccessMock.When(ctxCorrectUser, groupId).Then(nil)
	groupAccessChecker.CheckReadGroupAccessMock.When(ctxNotCorrectUser, groupId).Then(entity.ErrGroupWriteAccessDenied)
	cardService := NewCardsService(
		CardsServiceDeps{
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
	groupId := entity.GroupId(200)
	notExistGroupId := entity.GroupId(1)
	correctCards := []entity.Card{
		{},
	}
	tests := []struct {
		name    string
		ctx     context.Context
		groupId entity.GroupId
		result  []entity.Card
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
			err:     entity.ErrGroupWriteAccessDenied,
		},
		{
			name:    "no found",
			ctx:     ctx,
			groupId: notExistGroupId,
			err:     entity.ErrGroupNotFound,
		},
	}

	mc := minimock.NewController(t)

	cardsReaderMock := mocks.NewCardReaderMock(mc)
	cardsReaderMock.ListCardsMock.When(minimock.AnyContext, groupId).Then(correctCards, nil)
	cardsReaderMock.ListCardsMock.When(minimock.AnyContext, notExistGroupId).Then(nil, entity.ErrGroupNotFound)

	groupAccessChecker := mocks.NewGroupAccessChecker(mc)
	groupAccessChecker.CheckReadGroupAccessMock.When(ctx, groupId).Then(nil)
	groupAccessChecker.CheckReadGroupAccessMock.When(ctx, notExistGroupId).Then(nil)
	groupAccessChecker.CheckReadGroupAccessMock.When(
		ctxNotCorrectUser, groupId,
	).Then(entity.ErrGroupWriteAccessDenied)
	cardService := NewCardsService(
		CardsServiceDeps{
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
