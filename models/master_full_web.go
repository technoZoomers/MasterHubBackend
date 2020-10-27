package models

//easyjson:json
type MasterFull struct {
	UserId          int64    `json:"user_id"`
	Email    string    `json:"email"`
	Password string    `json:"password,omitempty"`
	Username        string   `json:"username"`
	Fullname        string   `json:"fullname"`
	Languages       []string `json:"language"`
	Theme           Theme    `json:"theme"`
	Description     string   `json:"description"`
	Qualification   string   `json:"qualification"`
	EducationFormat []string `json:"education_format"`
	AveragePrice    Price    `json:"avg_price"`
}
