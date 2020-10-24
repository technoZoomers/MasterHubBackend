package models

import "time"

type MessageDB struct {
	Id          int64
	UserId    int64
	ChatId int64
	Text        string
	Created    time.Time
}

type MessageInfoDB struct {
	Id          int64
	ChatId int64
	Text        string
	Created    time.Time
}