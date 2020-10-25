package models

import "time"

type MessageDB struct {
	Id          int64
	Info bool
	UserId    int64
	ChatId int64
	Text        string
	Created    time.Time
}
