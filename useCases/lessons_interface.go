package useCases

import "github.com/technoZoomers/MasterHubBackend/models"

type LessonsUCInterface interface {
	GetMastersLessons(masterId int64) (models.Lessons, error)
	CreateLesson(lesson *models.Lesson, masterId int64) error
	ChangeLessonInfo(lesson *models.Lesson, masterId int64, lessonId int64) error
	GetMastersLessonsRequests() (models.LessonRequests, error)
	CreateLessonRequest(lessonRequest *models.LessonRequest) error
	DeleteLessonRequest(studentId int64, lessonId int64) error
	DeleteMasterLesson(masterId int64, lessonId int64) error
	GetMastersLessonsStudents(masterId int64, lessonId int64) (models.LessonStudents, error)
}
