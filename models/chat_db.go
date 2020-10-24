package models

import "time"

type ChatDB struct {
	Id          int64
	Type        int64
	MasterId        int64
	StudentId        int64
	Created    time.Time
}
