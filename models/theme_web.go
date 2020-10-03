package models

//easyjson:json
type Theme struct {
	Id int64 `json:"id"`
	Theme string `json:"theme"`
	Subthemes []string `json:"subthemes"`
}


//easyjson:json
type Themes []Theme
