package models

//easyjson:json
type Master struct {
	UserId int64 `json:"userId"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Languages []string `json:"language"`
	Theme Theme `json:"theme"`
	Description string `json:"description"`
	Qualification string `json:"qualification"`
	EducationFormat []string `json:"educationFormat"`
	AveragePrice Price `json:"avgPrice"`
}

//easyjson:json
type Masters []Master