package model

import (
	"time"
)

const (
	CARD_SIDE_FIRST = CardSideNumber(iota)
	CARD_SIDE_SECOND
)

type CardId int64

type CardSideNumber int8

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

func (c *Card) SetText(first, second string) {
	if len(c.Sides) < 2 {
		c.Sides = make([]CardSide, 2)
	}
	c.Sides[CARD_SIDE_FIRST].Text = first
	c.Sides[CARD_SIDE_SECOND].Text = second
}

func (c *Card) GetFirst() CardSide {
	return c.GetText(CARD_SIDE_FIRST)
}
func (c *Card) GetSecond() CardSide {
	return c.GetText(CARD_SIDE_SECOND)
}
func (c *Card) GetText(side CardSideNumber) CardSide {
	return c.Sides[side]
}

type UpdateCard struct {
	Id        CardId
	FrontText string
	BackText  string
}
