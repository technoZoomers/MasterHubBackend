package models

//easyjson:json
type Master struct {
	UserId          int64    `json:"user_id"`
	Username        string   `json:"username"`
	Fullname        string   `json:"fullname"`
	Languages       []string `json:"language"`
	Theme           Theme    `json:"theme"`
	Description     string   `json:"description"`
	Qualification   string   `json:"qualification"`
	EducationFormat []string `json:"education_format"`
	AveragePrice    Price    `json:"hour_price"`
}

//easyjson:json
type Masters []Master
