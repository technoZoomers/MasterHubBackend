package models

type ChatsQueryValues struct {
	Type   int64
	Offset int64
	Limit  int64
}

type ChatsQueryValuesDB struct {
	Type   int64
	UserId int64
	User   int64 // master - 1, student - 2
	Limit  int64
	Offset int64
}
