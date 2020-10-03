package models

type ThemeDB struct {
	Id int64
	Name string
}

type SubthemeDB struct {
	Id int64
	ThemeId int64
	Name string
}
