package entity

import (
	"time"
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
	CardsProgress  []CardProgress
}

type CardProgress struct {
	Id                CardId
	LastReviewTime    time.Time
	AverageReviewTime time.Duration
	FailsCount        int
	HardCount         int
	GoodCount         int
	EasyCount         int
}
