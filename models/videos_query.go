package models

type VideosQueryValues struct {
	Subtheme []string
	Theme    string
	Popular  bool
	Old      bool
	Offset   int64
	Limit    int64
}

type VideosQueryValuesDB struct {
	Subtheme []int64
	Theme    []int64
	Popular  bool
	Old      bool
	Limit    int64
	Offset   int64
}
