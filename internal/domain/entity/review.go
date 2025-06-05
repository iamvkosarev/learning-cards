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

type Mark int32

const (
	MARK_NULL = Mark(iota)
	MARK_A
	MARK_B
	MARK_C
	MARK_D
	MARK_E
)

type GroupProgress struct {
	UserId         UserId
	GroupId        GroupId
	LastReviewTime time.Time
	CardsProgress  []CardReview
}

type CardReview struct {
	UserId   UserId
	GroupId  GroupId
	CardId   CardId
	Time     time.Time
	Duration time.Duration
	Answer   Answer
}

type CardMark struct {
	Id CardId
	Mark
}

type ReviewSettings struct {
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
