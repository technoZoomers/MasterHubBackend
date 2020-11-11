package repository

import (
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
)

type LessonsRepo struct {
	repository *Repository
}

func (lessonsRepo *LessonsRepo) InsertLesson(lesson *models.LessonDB) error {
	var dbError error
	transaction, err := lessonsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("INSERT INTO lessons (master_id, time_start, time_end, date, price, education_format, status) VALUES ($1, $2, $3, $4, $5, $6, $7) returning id",
		lesson.MasterId, lesson.TimeStart, lesson.TimeEnd, lesson.Date, lesson.Price, lesson.EducationFormat, lesson.Status)
	err = row.Scan(&lesson.Id)
	if err != nil {
		dbError = fmt.Errorf("failed to insert lesson: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := lessonsRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = lessonsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (lessonsRepo *LessonsRepo) GetMastersLessons() ([]models.LessonDB, error) {
	panic("implement me")
}

func (lessonsRepo *LessonsRepo) GetMastersLessonRequests() ([]models.LessonDB, error) {
	panic("implement me")
}
