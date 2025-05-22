package usecase

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/domain/contracts"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"sort"
	"time"
)

type ReviewUseCaseDeps struct {
	ProgressReader contracts.ProgressReader
	ProgressWriter contracts.ProgressWriter
	CardReader     contracts.CardReader
	GroupReader    contracts.GroupReader
}

type ReviewUseCase struct {
	ReviewUseCaseDeps
}

func NewReviewUseCase(deps ReviewUseCaseDeps) *ReviewUseCase {
	return &ReviewUseCase{deps}
}

func (r *ReviewUseCase) GetReviewCards(
	ctx context.Context, userId entity.UserId,
	groupId entity.GroupId, settings entity.ReviewSettings,
) (
	[]entity.Card,
	error,
) {
	op := "usecase.ReviewUseCase.GetReviewCards"

	group, err := getGroupAndCheckAccess(ctx, userId, groupId, op, r.GroupReader)
	if err != nil {
		return nil, err
	}

	cardsProgressRow, err := r.ProgressReader.GetCardsProgress(ctx, userId, groupId)
	if err != nil {
		return nil, fmt.Errorf("%s: error getting card progress: %w", op, err)
	}
	cardsProgress := make(map[entity.CardId]entity.CardProgress)
	for _, card := range cardsProgressRow {
		cardsProgress[card.Id] = card
	}
	cardsRow, err := r.CardReader.ListCards(ctx, groupId)
	if err != nil {
		return nil, fmt.Errorf("%s: error getting list of cards: %w", op, err)
	}
	cards := make(map[entity.CardId]entity.Card)
	for _, card := range cardsRow {
		cards[card.Id] = card
	}
	if len(cards) > len(cardsProgress) {
		for _, card := range cards {
			if _, ok := cardsProgress[card.Id]; !ok {
				cardsProgress[card.Id] = entity.CardProgress{
					Id:             card.Id,
					LastReviewTime: group.CreateTime,
				}
			}
		}
	}
	reviewCards := make([]entity.Card, 0)
	usedCards := make(map[entity.CardId]struct{})
	// AddCard new cards
	for cardId, cardProgress := range cardsProgress {
		if group.CreateTime.Equal(cardProgress.LastReviewTime) {
			reviewCards = append(reviewCards, cards[cardId])
			usedCards[cardId] = struct{}{}

			if len(reviewCards) >= settings.CardsCount {
				return reviewCards, nil
			}
		}
	}
	// AddCard long time no reviewed cards
	longTimeDuration := time.Hour * 24 * 3
	for cardId, cardProgress := range cardsProgress {
		if _, ok := usedCards[cardId]; ok {
			continue
		}
		if time.Now().After(cardProgress.LastReviewTime.Add(longTimeDuration)) {
			reviewCards = append(reviewCards, cards[cardId])
			usedCards[cardId] = struct{}{}

			if len(reviewCards) >= settings.CardsCount {
				return reviewCards, nil
			}
		}
	}

	// AddCard sorted by marks cards
	sortedByProgressCards := getSortedByProgressCards(cardsProgress, usedCards)
	for _, cardId := range sortedByProgressCards {
		reviewCards = append(reviewCards, cards[cardId])
		usedCards[cardId] = struct{}{}

		if len(reviewCards) >= settings.CardsCount {
			return reviewCards, nil
		}
	}

	return reviewCards, nil
}

func getSortedByProgressCards(
	progress map[entity.CardId]entity.CardProgress,
	used map[entity.CardId]struct{},
) []entity.CardId {
	cards := make([]entity.CardId, 0, len(progress)-len(used))
	for id := range progress {
		if _, ok := used[id]; ok {
			continue
		}
		cards = append(cards, id)
	}
	marks := make(map[entity.CardId]float64)
	for _, cardId := range cards {
		marks[cardId] = getCardMark(progress[cardId])
	}
	sort.Slice(
		cards, func(i, j int) bool {
			return marks[cards[i]] < marks[cards[j]]
		},
	)
	return cards
}

func getCardReviewsCount(progress entity.CardProgress) int {
	return progress.HardCount + progress.GoodCount + progress.FailsCount + progress.EasyCount
}

func getCardMark(progress entity.CardProgress) float64 {
	reviewsCount := getCardReviewsCount(progress)
	reviewsAvgValue := progress.FailsCount*1 + progress.HardCount*2 + progress.GoodCount*3 + progress.
		EasyCount*4
	return float64(reviewsAvgValue) / float64(reviewsCount)
}

func (r *ReviewUseCase) AddReviewResults(
	ctx context.Context, userId entity.UserId,
	groupId entity.GroupId, answers []entity.ReviewCardResult,
) error {
	op := "usecase.ReviewUseCase.GetReviewCards"

	_, err := getGroupAndCheckAccess(ctx, userId, groupId, op, r.GroupReader)
	if err != nil {
		return err
	}

	cardsProgressRow, err := r.ProgressReader.GetCardsProgress(ctx, userId, groupId)
	if err != nil {
		return err
	}
	cardsProgress := make(map[entity.CardId]entity.CardProgress)
	for _, card := range cardsProgressRow {
		cardsProgress[card.Id] = card
	}

	for _, answer := range answers {
		var cardPr entity.CardProgress
		var hasCardProgress bool
		if cardPr, hasCardProgress = cardsProgress[answer.CardId]; hasCardProgress {
			reviewsCount := int64(getCardReviewsCount(cardPr))
			cardPr.AverageReviewTime = time.Duration(
				(reviewsCount*cardPr.AverageReviewTime.
					Nanoseconds() + answer.Duration.Nanoseconds()) / (reviewsCount + 1),
			)
		} else {
			cardPr = entity.CardProgress{
				Id:                answer.CardId,
				AverageReviewTime: answer.Duration,
			}
		}
		switch answer.Answer {
		case entity.ANSWER_EASY:
			cardPr.EasyCount++
		case entity.ANSWER_GOOD:
			cardPr.GoodCount++
		case entity.ANSWER_HARD:
			cardPr.HardCount++
		case entity.ANSWER_FAIL:
			cardPr.FailsCount++
		}
		cardPr.LastReviewTime = time.Now()
		cardsProgress[answer.CardId] = cardPr
	}

	cardsProgressToSave := make([]entity.CardProgress, 0, len(cardsProgress))
	for _, progress := range cardsProgress {
		cardsProgressToSave = append(cardsProgressToSave, progress)
	}
	err = r.ProgressWriter.UpdateCardsProgress(ctx, userId, groupId, cardsProgressToSave)
	if err != nil {
		return err
	}
	return nil
}
