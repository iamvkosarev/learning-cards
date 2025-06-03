package service

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"sort"
	"time"
)

const (
	ANSWER_FAIL_SCORE = iota + 1.0
	ANSWER_HARD_CARDS
	ANSWER_GOOD_CARDS
	ANSWER_EASY_CARDS
)

const (
	MARK_A_START = float32(1 + (ANSWER_EASY_CARDS-iota)*(ANSWER_EASY_CARDS-1.0)/(ANSWER_EASY_CARDS+1.0))
	MARK_B_START
	MARK_C_START
	MARK_D_START
)

type ProgressReader interface {
	GetCardsProgress(ctx context.Context, user entity.UserId, group entity.GroupId) ([]entity.CardProgress, error)
}

type ProgressWriter interface {
	UpdateCardsProgress(
		ctx context.Context,
		user entity.UserId,
		group entity.GroupId,
		cardsProgress []entity.CardProgress,
	) error
}

type ReviewServiceDeps struct {
	ProgressReader ProgressReader
	ProgressWriter ProgressWriter
	CardReader     CardReader
	GroupReader    GroupReader
	Config         config.ReviewsService
}

type ReviewService struct {
	ReviewServiceDeps
}

func NewReviewService(deps ReviewServiceDeps) *ReviewService {
	return &ReviewService{deps}
}

func (r *ReviewService) GetReviewCards(
	ctx context.Context, userId entity.UserId,
	groupId entity.GroupId, settings entity.ReviewSettings,
) (
	[]entity.Card,
	error,
) {
	cards, progress, err := r.getCardsAndProgress(ctx, userId, groupId)
	if err != nil {
		return nil, err
	}

	reviewCards := make([]entity.Card, 0)
	usedCards := make(map[entity.CardId]struct{})
	// AddCard new cards
	for cardId, pr := range progress {
		if getCardReviewsCount(pr) == 0 {
			reviewCards = append(reviewCards, cards[cardId])
			usedCards[cardId] = struct{}{}

			if len(reviewCards) >= settings.CardsCount {
				return reviewCards, nil
			}
		}
	}
	// AddCard long time no reviewed cards
	longTimeDuration := time.Hour * 24 * 3
	for cardId, pr := range progress {
		if _, ok := usedCards[cardId]; ok {
			continue
		}
		if time.Now().After(pr.LastReviewTime.Add(longTimeDuration)) {
			reviewCards = append(reviewCards, cards[cardId])
			usedCards[cardId] = struct{}{}

			if len(reviewCards) >= settings.CardsCount {
				return reviewCards, nil
			}
		}
	}

	// AddCard sorted by marks cards
	sortedByProgressCards := r.getSortedCardsByScores(removeUniqueCards(progress, usedCards))
	for _, cardId := range sortedByProgressCards {
		reviewCards = append(reviewCards, cards[cardId])
		usedCards[cardId] = struct{}{}

		if len(reviewCards) >= settings.CardsCount {
			return reviewCards, nil
		}
	}

	return reviewCards, nil
}

func (r *ReviewService) GetCardsMarks(
	ctx context.Context,
	userId entity.UserId,
	groupId entity.GroupId,
) ([]entity.CardMark, error) {
	_, progress, err := r.getCardsAndProgress(ctx, userId, groupId)
	if err != nil {
		return nil, err
	}

	return r.getMarks(progress), nil
}

func (r *ReviewService) getCardsAndProgress(
	ctx context.Context,
	userId entity.UserId, groupId entity.GroupId,
) (
	map[entity.CardId]entity.
		Card, map[entity.CardId]entity.CardProgress, error,
) {
	cardsProgressRow, err := r.ProgressReader.GetCardsProgress(ctx, userId, groupId)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting card progress: %w", err)
	}
	cardsProgress := make(map[entity.CardId]entity.CardProgress)
	for _, card := range cardsProgressRow {
		cardsProgress[card.Id] = card
	}
	cardsRow, err := r.CardReader.ListCards(ctx, groupId)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting list of cards: %w", err)
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
					LastReviewTime: time.Now(),
				}
			}
		}
	}
	return cards, cardsProgress, nil
}

func removeUniqueCards[TCards any](
	cards map[entity.CardId]TCards,
	uniqueCards map[entity.CardId]struct{},
) map[entity.CardId]TCards {
	for id := range uniqueCards {
		delete(cards, id)
	}
	return cards
}

func (r *ReviewService) getSortedCardsByScores(
	progress map[entity.CardId]entity.CardProgress,
) []entity.CardId {
	cards := make([]entity.CardId, 0, len(progress))
	for id := range progress {
		cards = append(cards, id)
	}
	marks := r.getMarks(progress)
	marksMap := make(map[entity.CardId]entity.CardMark)
	for _, mark := range marks {
		marksMap[mark.Id] = mark
	}
	sort.Slice(
		cards, func(i, j int) bool {
			return marksMap[cards[i]].Mark < marksMap[cards[j]].Mark
		},
	)
	return cards
}

func (r *ReviewService) getMarks(progress map[entity.CardId]entity.CardProgress) []entity.CardMark {
	marks := make([]entity.CardMark, 0, len(progress))
	minAnswerDuration, maxAnswerDuration := getMinMaxAnswerDuration(progress)
	for id, pr := range progress {
		var mark entity.Mark
		reviewsCount := getCardReviewsCount(pr)
		if reviewsCount > 0 {
			answerScore := getCardAnswerScore(pr)
			durationScore := getDurationScore(
				minAnswerDuration.Seconds(), maxAnswerDuration.Seconds(),
				pr.AverageReviewTime.Seconds(),
			)
			score := float32(answerScore)*r.Config.AnswerInfluencePercent + float32(durationScore)*r.Config.SelectDurationInfluencePercent
			switch {
			case score > MARK_A_START:
				mark = entity.MARK_A
			case score > MARK_B_START:
				mark = entity.MARK_B
			case score > MARK_C_START:
				mark = entity.MARK_C
			case score > MARK_D_START:
				mark = entity.MARK_D
			default:
				mark = entity.MARK_E
			}
		} else {
			mark = entity.MARK_NULL
		}

		marks = append(
			marks, entity.CardMark{
				Mark: mark,
				Id:   id,
			},
		)
	}
	return marks
}

func getDurationScore(
	min float64,
	max float64,
	duration float64,
) float64 {
	if max-min <= 0 {
		return ANSWER_FAIL_SCORE
	}
	return (duration-min)/(max-min)*(ANSWER_EASY_CARDS-ANSWER_FAIL_SCORE) + ANSWER_FAIL_SCORE
}

func getMinMaxAnswerDuration(progress map[entity.CardId]entity.CardProgress) (
	min time.Duration,
	max time.Duration,
) {
	max = time.Duration(-1 << 63)
	min = time.Duration(int64(1<<63 - 1))
	for _, pr := range progress {
		if pr.AverageReviewTime > max {
			max = pr.AverageReviewTime
		}
		if pr.AverageReviewTime < min {
			min = pr.AverageReviewTime
		}
	}
	return min, max
}

func getCardAnswerScore(progress entity.CardProgress) float64 {
	reviewsCount := getCardReviewsCount(progress)
	reviewsAvgValue := progress.FailsCount*ANSWER_FAIL_SCORE +
		progress.HardCount*ANSWER_HARD_CARDS +
		progress.GoodCount*ANSWER_GOOD_CARDS +
		progress.EasyCount*ANSWER_EASY_CARDS
	return float64(reviewsAvgValue) / float64(reviewsCount)
}

func getCardReviewsCount(progress entity.CardProgress) int {
	return progress.HardCount + progress.GoodCount + progress.FailsCount + progress.EasyCount
}

func (r *ReviewService) AddReviewResults(
	ctx context.Context, userId entity.UserId,
	groupId entity.GroupId, answers []entity.ReviewCardResult,
) error {
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
