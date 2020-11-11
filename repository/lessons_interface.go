package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type LessonsRepoI interface {
	InsertLesson(lesson *models.LessonDB) error
	GetMastersLessons(masterId int64) ([]models.LessonDB, error)
	GetMastersLessonRequests() ([]models.LessonDB, error)
	GetLessonByIdAndMasterId(lesson *models.LessonDB, lessonId int64, masterId int64) error
	CheckLessonTimeRange(lesson *models.LessonDB) ([]int64, error)
	UpdateLessonByIdAndMasterId(lesson *models.LessonDB) error
}
