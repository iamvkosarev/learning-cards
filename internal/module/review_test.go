package module

import (
	"fmt"
	"github.com/gojuno/minimock/v3"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/model"
	"github.com/iamvkosarev/learning-cards/internal/module/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

func TestReviewService_DeleteNotUsedReviews(t *testing.T) {
	groupId := model.GroupId(200)
	userId := model.UserId(200)
	tests := []struct {
		name     string
		cards    []*model.Card
		reviews  []*model.CardReview
		result   map[model.CardId]int
		settings model.ReviewSettings
		config   config.ReviewsService
		err      error
	}{
		{
			name: "success_get_old_cards",
			cards: []*model.Card{
				{Id: model.CardId(1)},
				{Id: model.CardId(2)},
			},
			reviews: []*model.CardReview{
				{
					CardId: model.CardId(1),
					Answer: model.ANSWER_FAIL,
					Time:   time.Now(),
				},
				{
					CardId: model.CardId(2),
					Answer: model.ANSWER_EASY,
					Time:   time.Now().Add(-model.NeedToDoReviewDuration).Add(-1 * time.Hour),
				},
			},
			result: map[model.CardId]int{
				model.CardId(2): 0,
			},
			settings: model.ReviewSettings{
				CardsCount: 1,
			},
			config: config.ReviewsService{
				AnswerInfluencePercent: 1,
				ReviewStepWeight:       1,
			},
			err: nil,
		},
		{
			name: "success_get_new_card",
			cards: []*model.Card{
				{Id: model.CardId(1)},
				{Id: model.CardId(2)},
			},
			reviews: []*model.CardReview{
				{
					CardId: model.CardId(1),
					Answer: model.ANSWER_FAIL,
					Time:   time.Now(),
				},
			},
			result: map[model.CardId]int{
				model.CardId(2): 0,
			},
			settings: model.ReviewSettings{
				CardsCount: 1,
			},
			config: config.ReviewsService{
				AnswerInfluencePercent: 1,
				ReviewStepWeight:       1,
			},
			err: nil,
		},
		{
			name: "success_get_failed_card",
			cards: []*model.Card{
				{Id: model.CardId(1)},
				{Id: model.CardId(2)},
			},
			reviews: []*model.CardReview{
				{
					CardId: model.CardId(1),
					Answer: model.ANSWER_GOOD,
					Time:   time.Now(),
				},
				{
					CardId: model.CardId(2),
					Answer: model.ANSWER_FAIL,
					Time:   time.Now(),
				},
			},
			result: map[model.CardId]int{
				model.CardId(2): 0,
			},
			settings: model.ReviewSettings{
				CardsCount: 1,
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
				cards, err := reviewsService.GetReviewCards(minimock.AnyContext, groupId, test.settings)
				if test.err == nil {
					if len(test.result) != len(cards) {
						assert.Fail(
							t, fmt.Sprintf("result length mismatch: expected %d, got %d", len(test.result), len(cards)),
						)
					} else {
						for _, card := range cards {
							if _, ok := test.result[card.Id]; !ok {
								assert.Fail(t, fmt.Sprintf("card %d not found in result", card.Id))
							} else {
								test.result[card.Id]++
								if test.result[card.Id] > 1 {
									assert.Fail(t, fmt.Sprintf("card %d returned twice in result", card.Id))
								}
							}
						}
						for id, count := range test.result {
							if count != 1 {
								assert.Fail(t, fmt.Sprintf("expected card with id: %d", id))
							}
						}
					}
				}
				assert.Equal(t, test.err, err)
			},
		)
	}
}
