package models

import "time"

//easyjson:json
type Message struct {
	Id          int64    `json:"id"`
	Type        int64   `json:"type"`
	AuthorId    int64   `json:"author_id,omitempty"`
	ChatId int64 `json:"chat_id"`
	Text        string   `json:"text"`
	Created    time.Time `json:"created"`
}

//easyjson:json
type Messages []Message
