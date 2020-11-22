package models

//easyjson:json
type Notification struct {
	Text   string `json:"text"`
	UserId int64  `json:"user_id"`
}
