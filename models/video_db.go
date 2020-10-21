package models

import "time"

type VideoDB struct {
	Id          int64
	MasterId    int64
	Filename    string
	Extension   string
	Name        string
	Description string
	Intro       bool
	Theme       int64
	Rating int64
	Uploaded    time.Time
}
