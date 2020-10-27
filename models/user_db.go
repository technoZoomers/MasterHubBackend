package models

import "time"

type UserDB struct {
	Id       int64
	Email    string
	Password string
	Type     int64
	Created  time.Time
}
