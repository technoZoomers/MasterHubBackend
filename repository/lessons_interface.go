package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type LessonsRepoI interface {
	InsertLesson(lesson *models.LessonDB) error
	InsertLessonRequest(lessonRequest *models.LessonStudentDB) error
	GetLessonRequestByStudentIdAndLessonId(lessonRequest *models.LessonStudentDB, studentId int64, lessonId int64) error
	GetLessonRequestByStudentUserIdAndLessonId(lessonRequest *models.LessonStudentDB, studentId int64, lessonId int64) error
	DeleteLessonRequestByStudentIdAndLessonId(studentId int64, lessonId int64) error
	GetMastersLessons(masterId int64) ([]models.LessonDB, error)
	GetLessonByIdAndMasterId(lesson *models.LessonDB, lessonId int64, masterId int64) error
	GetLessonById(lesson *models.LessonDB, lessonId int64) error
	CheckLessonTimeRange(lesson *models.LessonDB) ([]int64, error)
	UpdateLessonByIdAndMasterId(lesson *models.LessonDB) error
	DeleteLessonById(lessonId int64) error
	GetLessonStudents(lessonId int64) ([]int64, error)
	GetMastersLessonsRequests(masterId int64) ([]models.LessonStudentDB, error)
	UpdateLessonRequest(lessonRequest *models.LessonStudentDB) error
}
