package model

import (
	"time"
)

const (
	GROUP_VISIBILITY_NULL = GroupVisibility(iota)
	GROUP_VISIBILITY_PRIVATE
	GROUP_VISIBILITY_PUBLIC
	GROUP_VISIBILITY_UNLISTED
)

const (
	CARD_SIDE_TYPE_NULL = CardSideType(iota)
	CARD_SIDE_TYPE_JAPANESE
)

type CardSideType int8

type GroupVisibility uint8
type GroupId int64

type Group struct {
	Id          GroupId
	OwnerId     UserId
	Name        string
	Description string
	CreateTime  time.Time
	UpdateTime  time.Time
	Visibility  GroupVisibility
	// SideTypes is list of cards side's types.
	// It could be appointed if it equals to CARD_SIDE_TYPE_NULL
	// and used to decorate cards.
	CardSideTypes []CardSideType
}

type UpdateGroup struct {
	Id           GroupId
	Name         string
	Description  string
	Visibility   GroupVisibility
	CardSideType []CardSideType
}
