package entity

import (
	"time"
)

const (
	GROUP_VISIBILITY_NULL = GroupVisibility(iota)
	GROUP_VISIBILITY_PRIVATE
	GROUP_VISIBILITY_PUBLIC
	GROUP_VISIBILITY_UNLISTED
)

type GroupVisibility int32
type GroupId int64

type Group struct {
	Id          GroupId
	OwnerId     UserId
	Name        string
	Description string
	CreateTime  time.Time
	Visibility  GroupVisibility
}
