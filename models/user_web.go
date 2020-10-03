package models

import "time"

//easyjson:json
type User struct {
	Id int64 `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
	Type int64 `json:"type"`
	Created  time.Time `json:"created"`
}

//easyjson:json
type Users []User