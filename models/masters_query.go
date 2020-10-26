package models

type MastersQueryValues struct {
	Subtheme        []string
	Theme           string
	Qualification   string
	EducationFormat string
	Language        []string
	Search          string
	Offset          int64
	Limit           int64
}

type MastersQueryValuesDB struct {
	Subtheme        []int64
	Theme           []int64
	Qualification   int64
	EducationFormat int64
	Language        []int64
	Limit           int64
	Offset          int64
}
