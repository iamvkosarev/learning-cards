package model

import (
	"time"
)

const (
	SIDE_FIRST = iota
	SIDE_SECOND
)

type CardId int64

type Decoration interface {
	ToString() string
}

type ReadingPair struct {
	Text    string
	Reading string
}

type CardSide struct {
	Text         string
	ReadingPairs []ReadingPair
}

type Card struct {
	Id         CardId
	GroupId    GroupId
	Sides      []CardSide
	CreateTime time.Time
	UpdateTime time.Time
}

func NewCardSides(front, back string) []CardSide {
	return []CardSide{
		{Text: front},
		{Text: back},
	}
}

type UpdateCard struct {
	Id        CardId
	SidesText []string
}
