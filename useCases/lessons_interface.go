package useCases

import "github.com/technoZoomers/MasterHubBackend/models"

type LessonsUCInterface interface {
	GetMastersLessons(masterId int64) (models.Lessons, error)
	CreateLesson(lesson *models.Lesson, masterId int64) error
	ChangeLessonInfo(lesson *models.Lesson, masterId int64, lessonId int64) error
	GetMastersLessonsRequests() (models.LessonRequests, error)
	CreateLessonRequest(studentId int64, lessonId int64) error
	DeleteLessonRequest(lessonRequest *models.LessonRequest) error
}