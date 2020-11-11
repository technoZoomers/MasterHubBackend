package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type LessonDB struct {
	Id              int64
	MasterId        int64
	TimeStart       string
	TimeEnd         string
	Date            time.Time
	EducationFormat int64
	Price           decimal.Decimal
	Status          int64
}
