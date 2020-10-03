package models

import "time"

//easyjson:json
type VideoData struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Description string `json:"description,omitempty"`
	Theme Theme `json:"theme,omitempty"`
	Uploaded  time.Time `json:"uploaded"`
}

//easyjson:json
type VideosData []VideoData
