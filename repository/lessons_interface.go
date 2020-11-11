package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type LessonsRepoI interface {
	InsertLesson(lesson *models.LessonDB) error
	GetMastersLessons() ([]models.LessonDB, error)
	GetMastersLessonRequests() ([]models.LessonDB, error)
}
