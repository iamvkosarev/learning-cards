package entity

import (
	"time"
)

type Card struct {
	Id         CardId
	GroupId    GroupId
	FrontText  string
	BackText   string
	CreateTime time.Time
}

type CardId int64
