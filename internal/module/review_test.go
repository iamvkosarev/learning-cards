package module

import (
	"github.com/gojuno/minimock/v3"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/model"
	"github.com/iamvkosarev/learning-cards/internal/module/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReviewService_GetCardsMarks(t *testing.T) {
	groupId := model.GroupId(200)
	userId := model.UserId(200)
	tests := []struct {
		name    string
		cards   []*model.Card
		reviews []*model.CardReview
		result  []model.CardMark
		config  config.ReviewsService
		err     error
	}{
		{
			name:    "success",
			cards:   []*model.Card{},
			reviews: []*model.CardReview{},
			result:  []model.CardMark{},
			err:     nil,
		},
		{
			name: "null_reviews",
			cards: []*model.Card{
				{},
			},
			reviews: []*model.CardReview{},
			result: []model.CardMark{
				{
					Mark: model.MARK_NULL,
				},
			},
			err: nil,
		},
		{
			name: "fail_answer_mark",
			cards: []*model.Card{
				{},
			},
			reviews: []*model.CardReview{
				{Answer: model.ANSWER_FAIL},
			},
			result: []model.CardMark{
				{
					Mark: model.MARK_E,
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
			cards: []*model.Card{
				{},
			},
			reviews: []*model.CardReview{
				{Answer: model.ANSWER_EASY},
			},
			result: []model.CardMark{
				{
					Mark: model.MARK_A,
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

				reviewsService := NewReviews(
					ReviewsDeps{
						UserVerifier: userVerifierMock,
						ReviewReader: reviewReaderMock,
						CardReader:   cardReaderMock,
						Config:       test.config,
					},
				)

				cardsMarks, err := reviewsService.GetCardsMarks(
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
