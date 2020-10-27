package models

//easyjson:json
type StudentFull struct {
	UserId    int64    `json:"user_id"`
	Email    string    `json:"email"`
	Password string    `json:"password,omitempty"`
	Username  string   `json:"username"`
	Fullname  string   `json:"fullname"`
	Languages []string `json:"language"`
}
