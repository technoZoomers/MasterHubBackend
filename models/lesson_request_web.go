package models

//easyjson:json
type LessonRequest struct {
	LessonId  int64 `json:"lesson_id"`
	StudentId int64 `json:"student_id"`
	Status    int64 `json:"status"`
}

//easyjson:json
type LessonRequests []LessonRequest
