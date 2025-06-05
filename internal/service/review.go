package service

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"slices"
	"sort"
	"time"
)

const (
	ANSWER_FAIL_SCORE = iota + 1.0
	ANSWER_HARD_SCORE
	ANSWER_GOOD_SCORE
	ANSWER_EASY_SCORE
)

const (
	MARK_A_START = float32(1 + (ANSWER_EASY_SCORE-iota)*(ANSWER_EASY_SCORE-1.0)/(ANSWER_EASY_SCORE+1.0))
	MARK_B_START
	MARK_C_START
	MARK_D_START
)

//go:generate minimock -i ReviewReader -o ./mocks/review_reader_mock.go -n ReviewReaderMock -p mocks
type ReviewReader interface {
	GetCardsReviews(ctx context.Context, user entity.UserId, group entity.GroupId) ([]entity.CardReview, error)
}

//go:generate minimock -i ReviewWriter -o ./mocks/review_writer_mock.go -n ReviewWriterMock -p mocks
type ReviewWriter interface {
	AddCardsReviews(
		ctx context.Context,
		user entity.UserId,
		group entity.GroupId,
		cardsProgress []entity.CardReview,
	) error
	DeleteNotUsedReviews(
		ctx context.Context, userId entity.UserId, groupId entity.GroupId,
	) error
}

type ReviewServiceDeps struct {
	ReviewReader ReviewReader
	ReviewWriter ReviewWriter
	UserVerifier UserVerifier
	CardReader   CardReader
	GroupReader  GroupReader
	Config       config.ReviewsService
}

type ReviewService struct {
	ReviewServiceDeps
}

func NewReviewService(deps ReviewServiceDeps) *ReviewService {
	return &ReviewService{deps}
}

func (r *ReviewService) GetReviewCards(
	ctx context.Context,
	groupId entity.GroupId, settings entity.ReviewSettings,
) (
	[]entity.Card,
	error,
) {
	userId, err := r.UserVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}

	cards, cardsReviews, err := r.getCardsAndReviews(ctx, userId, groupId)
	if err != nil {
		return nil, err
	}

	reviewCards := make([]entity.Card, 0)
	usedCards := make(map[entity.CardId]struct{})
	// AddCard new cards
	for cardId, cardReviews := range cardsReviews {
		if len(cardReviews) == 0 {
			reviewCards = append(reviewCards, cards[cardId])
			usedCards[cardId] = struct{}{}

			if len(reviewCards) >= settings.CardsCount {
				return reviewCards, nil
			}
		}
	}
	// AddCard long time no reviewed cards
	longTimeDuration := time.Hour * 24 * 3
	for cardId, cardReviews := range cardsReviews {
		if _, ok := usedCards[cardId]; ok || len(cardReviews) == 0 {
			continue
		}
		lastReviewTime := cardReviews[0].Time
		for _, card := range cardReviews {
			if card.Time.After(lastReviewTime) {
				lastReviewTime = card.Time
			}
		}
		if time.Now().After(lastReviewTime.Add(longTimeDuration)) {
			reviewCards = append(reviewCards, cards[cardId])
			usedCards[cardId] = struct{}{}

			if len(reviewCards) >= settings.CardsCount {
				return reviewCards, nil
			}
		}
	}

	// AddCard sorted by marks cards
	sortedByProgressCards := r.getSortedCardsByScores(removeUniqueCards(cardsReviews, usedCards))
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
	groupId entity.GroupId,
) ([]entity.CardMark, error) {
	userId, err := r.UserVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return nil, err
	}

	_, progress, err := r.getCardsAndReviews(ctx, userId, groupId)
	if err != nil {
		return nil, err
	}

	return r.getMarks(progress), nil
}

func (r *ReviewService) getCardsAndReviews(
	ctx context.Context,
	userId entity.UserId, groupId entity.GroupId,
) (
	map[entity.CardId]entity.
		Card, map[entity.CardId][]entity.CardReview, error,
) {
	cardsProgressRow, err := r.ReviewReader.GetCardsReviews(ctx, userId, groupId)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting card progress: %w", err)
	}
	cardsProgress := make(map[entity.CardId][]entity.CardReview)
	for _, pr := range cardsProgressRow {
		if _, ok := cardsProgress[pr.CardId]; !ok {
			cardsProgress[pr.CardId] = make([]entity.CardReview, 0)
		}
		cardsProgress[pr.CardId] = append(cardsProgress[pr.CardId], pr)
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
				cardsProgress[card.Id] = make([]entity.CardReview, 0)
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
	reviews map[entity.CardId][]entity.CardReview,
) []entity.CardId {
	cards := make([]entity.CardId, 0, len(reviews))
	for id := range reviews {
		cards = append(cards, id)
	}
	marks := r.getMarks(reviews)
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

func (r *ReviewService) getMarks(cardsReviews map[entity.CardId][]entity.CardReview) []entity.CardMark {
	marks := make([]entity.CardMark, 0, len(cardsReviews))
	durationScores := getCardsDurationScores(cardsReviews)
	for id, cardReviews := range cardsReviews {
		var mark entity.Mark
		if len(cardReviews) > 0 {
			answerScore := getCardAnswerScore(r.Config.ReviewStepWeight, cardReviews)
			durationScore := durationScores[id]
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

func getCardsDurationScores(cardReviews map[entity.CardId][]entity.CardReview) map[entity.CardId]float64 {
	minArg := time.Duration(-1 << 63)
	maxArg := time.Duration(int64(1<<63 - 1))
	argDurations := make(map[entity.CardId]time.Duration)
	for id, reviews := range cardReviews {
		if len(reviews) == 0 {
			argDurations[id] = 0
			continue
		}
		sumDuration := time.Duration(0)
		for _, review := range reviews {
			sumDuration += review.Duration
		}
		argDuration := sumDuration / time.Duration(len(reviews))
		if sumDuration > maxArg {
			maxArg = sumDuration
		}
		if sumDuration < minArg {
			minArg = sumDuration
		}
		argDurations[id] = argDuration
	}
	scores := make(map[entity.CardId]float64)
	for id, reviews := range cardReviews {
		if len(reviews) == 0 {
			scores[id] = 0
			continue
		}
		argDuration := argDurations[id]
		score := (argDuration-minArg)/(maxArg-minArg)*(ANSWER_EASY_SCORE-ANSWER_FAIL_SCORE) + ANSWER_FAIL_SCORE
		scores[id] = float64(score)
	}
	return scores
}

func getCardAnswerScore(weightStep float64, reviews []entity.CardReview) float64 {
	scores := 0.0
	slices.SortFunc(
		reviews, func(a, b entity.CardReview) int {
			return a.Time.Compare(b.Time)
		},
	)
	weights := 0.0
	for i, review := range reviews {
		weight := weightStep * float64(1+i)
		weights += weight
		scores += float64(getAnswerScore(review.Answer)) * weight
	}
	return scores / weights
}

func getAnswerScore(answer entity.Answer) int {
	switch answer {
	case entity.ANSWER_EASY:
		return ANSWER_EASY_SCORE
	case entity.ANSWER_GOOD:
		return ANSWER_GOOD_SCORE
	case entity.ANSWER_HARD:
		return ANSWER_HARD_SCORE
	case entity.ANSWER_FAIL:
		return ANSWER_FAIL_SCORE
	}
	return 0
}

func (r *ReviewService) AddReviewResults(
	ctx context.Context,
	groupId entity.GroupId, answers []entity.ReviewCardResult,
) error {
	userId, err := r.UserVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return err
	}

	reviews := make([]entity.CardReview, 0, len(answers))
	reviewedCards := make([]entity.CardId, 0, len(answers))
	for _, answer := range answers {
		reviews = append(
			reviews, entity.CardReview{
				UserId:   userId,
				GroupId:  groupId,
				CardId:   answer.CardId,
				Time:     time.Now(),
				Duration: answer.Duration,
				Answer:   answer.Answer,
			},
		)
		reviewedCards = append(reviewedCards, answer.CardId)
	}

	if err := r.ReviewWriter.AddCardsReviews(ctx, userId, groupId, reviews); err != nil {
		return err
	}
	if err := r.ReviewWriter.DeleteNotUsedReviews(ctx, userId, groupId); err != nil {
		return err
	}

	return nil
}
