package models

//easyjson:json
type Status struct {
	Online    bool `json:"online"`
	IsCalling bool `json:"is_calling"`
}
