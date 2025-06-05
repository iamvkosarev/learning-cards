package service

import (
	"github.com/gojuno/minimock/v3"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"github.com/iamvkosarev/learning-cards/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReviewService_GetCardsMarks(t *testing.T) {
	groupId := entity.GroupId(200)
	userId := entity.UserId(200)
	tests := []struct {
		name    string
		cards   []entity.Card
		reviews []entity.CardReview
		result  []entity.CardMark
		config  config.ReviewsService
		err     error
	}{
		{
			name:    "success",
			cards:   []entity.Card{},
			reviews: []entity.CardReview{},
			result:  []entity.CardMark{},
			err:     nil,
		},
		{
			name: "null_reviews",
			cards: []entity.Card{
				{},
			},
			reviews: []entity.CardReview{},
			result: []entity.CardMark{
				{
					Mark: entity.MARK_NULL,
				},
			},
			err: nil,
		},
		{
			name: "fail_answer_mark",
			cards: []entity.Card{
				{},
			},
			reviews: []entity.CardReview{
				{Answer: entity.ANSWER_FAIL},
			},
			result: []entity.CardMark{
				{
					Mark: entity.MARK_E,
				},
			},
			config: config.ReviewsService{
				AnswerInfluencePercent: 1,
				ReviewStepWeight:       1,
			},
			err: nil,
		},
		{
			name: "fail_answer_mark",
			cards: []entity.Card{
				{},
			},
			reviews: []entity.CardReview{
				{Answer: entity.ANSWER_EASY},
			},
			result: []entity.CardMark{
				{
					Mark: entity.MARK_A,
				},
			},
			config: config.ReviewsService{
				AnswerInfluencePercent: 1,
				ReviewStepWeight:       1,
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				mc := minimock.NewController(t)

				userVerifierMock := mocks.NewNewUserVerifierMock(mc)
				userVerifierMock.VerifyUserByContextMock.Return(userId, nil)

				reviewReaderMock := mocks.NewReviewReaderMock(mc)
				reviewReaderMock.GetCardsReviewsMock.Return(test.reviews, nil)

				cardReaderMock := mocks.NewCardReaderMock(mc)
				cardReaderMock.ListCardsMock.Return(test.cards, nil)

				reviewService := NewReviewService(
					ReviewServiceDeps{
						UserVerifier: userVerifierMock,
						ReviewReader: reviewReaderMock,
						CardReader:   cardReaderMock,
						Config:       test.config,
					},
				)

				cardsMarks, err := reviewService.GetCardsMarks(
					minimock.AnyContext, groupId,
				)
				if test.err == nil {
					assert.Equal(t, test.result, cardsMarks)
				}
				assert.Equal(t, test.err, err)
			},
		)
	}
}
