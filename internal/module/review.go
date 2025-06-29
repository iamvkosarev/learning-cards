package module

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"github.com/iamvkosarev/learning-cards/internal/model"
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
	GetCardsReviews(ctx context.Context, user model.UserId, group model.GroupId) ([]*model.CardReview, error)
}

//go:generate minimock -i ReviewWriter -o ./mocks/review_writer_mock.go -n ReviewWriterMock -p mocks
type ReviewWriter interface {
	AddCardsReviews(
		ctx context.Context,
		user model.UserId,
		group model.GroupId,
		cardsProgress []model.CardReview,
	) error
	DeleteNotUsedReviews(
		ctx context.Context, userId model.UserId, groupId model.GroupId,
	) error
}

type ReviewsDeps struct {
	ReviewReader ReviewReader
	ReviewWriter ReviewWriter
	UserVerifier UserVerifier
	CardReader   CardReader
	GroupReader  GroupReader
	Config       config.ReviewsService
}

type Reviews struct {
	ReviewsDeps
}

func NewReviews(deps ReviewsDeps) *Reviews {
	return &Reviews{deps}
}

func (r *Reviews) GetReviewCards(
	ctx context.Context,
	groupId model.GroupId, settings model.ReviewSettings,
) (
	[]*model.Card,
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

	reviewCards := make([]*model.Card, 0)
	usedCards := make(map[model.CardId]struct{})
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
		if time.Now().After(lastReviewTime.Add(model.NeedToDoReviewDuration)) {
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

func (r *Reviews) GetCardsMarks(
	ctx context.Context,
	groupId model.GroupId,
) ([]model.CardMark, error) {
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

func (r *Reviews) getCardsAndReviews(
	ctx context.Context,
	userId model.UserId, groupId model.GroupId,
) (
	map[model.CardId]*model.
		Card, map[model.CardId][]*model.CardReview, error,
) {
	cardsProgressRow, err := r.ReviewReader.GetCardsReviews(ctx, userId, groupId)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting card progress: %w", err)
	}
	cardsProgress := make(map[model.CardId][]*model.CardReview)
	for _, pr := range cardsProgressRow {
		if _, ok := cardsProgress[pr.CardId]; !ok {
			cardsProgress[pr.CardId] = make([]*model.CardReview, 0)
		}
		cardsProgress[pr.CardId] = append(cardsProgress[pr.CardId], pr)
	}
	cardsRow, err := r.CardReader.ListCards(ctx, groupId)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting list of cards: %w", err)
	}
	cards := make(map[model.CardId]*model.Card)
	for _, card := range cardsRow {
		cards[card.Id] = card
	}
	if len(cards) > len(cardsProgress) {
		for _, card := range cards {
			if _, ok := cardsProgress[card.Id]; !ok {
				cardsProgress[card.Id] = make([]*model.CardReview, 0)
			}
		}
	}
	return cards, cardsProgress, nil
}

func removeUniqueCards[TCards any](
	cards map[model.CardId]TCards,
	uniqueCards map[model.CardId]struct{},
) map[model.CardId]TCards {
	for id := range uniqueCards {
		delete(cards, id)
	}
	return cards
}

func (r *Reviews) getSortedCardsByScores(
	reviews map[model.CardId][]*model.CardReview,
) []model.CardId {
	cards := make([]model.CardId, 0, len(reviews))
	for id := range reviews {
		cards = append(cards, id)
	}
	marks := r.getMarks(reviews)
	marksMap := make(map[model.CardId]model.CardMark)
	for _, mark := range marks {
		marksMap[mark.Id] = mark
	}
	sort.Slice(
		cards, func(i, j int) bool {
			return marksMap[cards[i]].Mark > marksMap[cards[j]].Mark
		},
	)
	return cards
}

func (r *Reviews) getMarks(cardsReviews map[model.CardId][]*model.CardReview) []model.CardMark {
	marks := make([]model.CardMark, 0, len(cardsReviews))
	durationScores := getCardsDurationScores(cardsReviews)
	for id, cardReviews := range cardsReviews {
		var mark model.Mark
		if len(cardReviews) > 0 {
			answerScore := getCardAnswerScore(r.Config.ReviewStepWeight, cardReviews)
			durationScore := durationScores[id]
			score := float32(answerScore)*r.Config.AnswerInfluencePercent + float32(durationScore)*r.Config.SelectDurationInfluencePercent
			switch {
			case score > MARK_A_START:
				mark = model.MARK_A
			case score > MARK_B_START:
				mark = model.MARK_B
			case score > MARK_C_START:
				mark = model.MARK_C
			case score > MARK_D_START:
				mark = model.MARK_D
			default:
				mark = model.MARK_E
			}
		} else {
			mark = model.MARK_NULL
		}

		marks = append(
			marks, model.CardMark{
				Mark: mark,
				Id:   id,
			},
		)
	}
	return marks
}

func getCardsDurationScores(cardReviews map[model.CardId][]*model.CardReview) map[model.CardId]float64 {
	minArg := time.Duration(-1 << 63)
	maxArg := time.Duration(int64(1<<63 - 1))
	argDurations := make(map[model.CardId]time.Duration)
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
	scores := make(map[model.CardId]float64)
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

func getCardAnswerScore(weightStep float64, reviews []*model.CardReview) float64 {
	scores := 0.0
	slices.SortFunc(
		reviews, func(a, b *model.CardReview) int {
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

func getAnswerScore(answer model.Answer) int {
	switch answer {
	case model.ANSWER_EASY:
		return ANSWER_EASY_SCORE
	case model.ANSWER_GOOD:
		return ANSWER_GOOD_SCORE
	case model.ANSWER_HARD:
		return ANSWER_HARD_SCORE
	case model.ANSWER_FAIL:
		return ANSWER_FAIL_SCORE
	}
	return 0
}

func (r *Reviews) AddReviewResults(
	ctx context.Context,
	groupId model.GroupId, answers []model.ReviewCardResult,
) error {
	userId, err := r.UserVerifier.VerifyUserByContext(ctx)
	if err != nil {
		return err
	}

	reviews := make([]model.CardReview, 0, len(answers))
	reviewedCards := make([]model.CardId, 0, len(answers))
	for _, answer := range answers {
		reviews = append(
			reviews, model.CardReview{
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
