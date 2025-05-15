package entity

import (
	"time"
)

type ReviewMode int32
type Answer int32

const (
	REVIEW_MODE_BASIC = ReviewMode(iota)
)

const (
	ANSWER_EASY = Answer(iota)
	ANSWER_GOOD
	ANSWER_HARD
	ANSWER_FAIL
)

type GroupReviewInfo struct {
	GroupId    GroupId
	CardsCount int
}

type UpdateGroupReviewInfo struct {
	GroupId    GroupId
	CardsCount int
}

type ReviewCardResult struct {
	CardId   CardId
	Answer   Answer
	Duration time.Duration
}
