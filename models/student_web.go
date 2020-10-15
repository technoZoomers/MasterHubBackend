package models

//easyjson:json
type Student struct {
	UserId    int64    `json:"userId"`
	Username  string   `json:"username"`
	Fullname  string   `json:"fullname"`
	Languages []string `json:"language"`
}

//easyjson:json
type Students []Student
