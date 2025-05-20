package entity

import (
	"time"
)

type CardId int64

type Card struct {
	Id         CardId
	GroupId    GroupId
	FrontText  string
	BackText   string
	CreateTime time.Time
	UpdateTime time.Time
}

type UpdateCard struct {
	Id        CardId
	FrontText string
	BackText  string
}
