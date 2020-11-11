package models

//easyjson:json
type Lesson struct {
	Id              int64  `json:"id"`
	MasterId        int64  `json:"master_id"`
	TimeStart       string `json:"time_start"`
	TimeEnd         string `json:"time_end"`
	Duration        string `json:"duration"`
	Date            string `json:"date"`
	EducationFormat string `json:"education_format"`
	Price           Price  `json:"price"`
	Status          int64  `json:"status"`
}

//easyjson:json
type Lessons []Lesson
