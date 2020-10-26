package models

import "time"

//easyjson:json
type Chat struct {
	Id        int64     `json:"id"`
	Type      int64     `json:"type"`
	MasterId  int64     `json:"master_id"`
	StudentId int64     `json:"student_id"`
	Created   time.Time `json:"created"`
}

//easyjson:json
type Chats []Chat
