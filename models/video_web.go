package models

import "time"

//easyjson:json
type VideoData struct {
	Id          int64     `json:"id,omitempty"`
	Intro 		bool      `json:"intro"`
	Name        string    `json:"name"`
	FileExt     string    `json:"extension"`
	Description string    `json:"description,omitempty"`
	Theme       Theme     `json:"theme,omitempty"`
	Uploaded    time.Time `json:"uploaded"`
}

//easyjson:json
type VideosData []VideoData
