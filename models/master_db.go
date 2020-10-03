package models

import "github.com/shopspring/decimal"

type MasterDB struct {
	Id int64
	UserId int64
	Username string
	Fullname string
	Theme int64
	Description string
	Qualification int64
	EducationFormat int64
	AveragePrice decimal.Decimal
}
