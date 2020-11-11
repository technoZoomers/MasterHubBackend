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

func (lessonsRepo *LessonsRepo) GetMastersLessons(masterId int64) ([]models.LessonDB, error) {
	var dbError error
	lessons := make([]models.LessonDB, 0)
	transaction, err := lessonsRepo.repository.startTransaction()
	if err != nil {
		return lessons, err
	}
	rows, err := transaction.Query(`SELECT * FROM lessons WHERE master_id = $1`, masterId)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve lessons: %v", err.Error())
		logger.Errorf(dbError.Error())
		return lessons, dbError
	}
	for rows.Next() {
		var lessFound models.LessonDB
		err = rows.Scan(&lessFound.Id, &lessFound.MasterId, &lessFound.TimeStart,
			&lessFound.TimeEnd, &lessFound.Date, &lessFound.Price, &lessFound.EducationFormat, &lessFound.Status)
		if err != nil {
			dbError = fmt.Errorf("failed to retrieve lesson: %v", err)
			logger.Errorf(dbError.Error())
			return lessons, dbError
		}
		lessons = append(lessons, lessFound)
	}
	err = lessonsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return lessons, err
	}
	return lessons, nil
}

func (lessonsRepo *LessonsRepo) GetMastersLessonRequests() ([]models.LessonDB, error) {
	panic("implement me")
}
func (lessonsRepo *LessonsRepo) CheckLessonTimeRange(lesson *models.LessonDB) ([]int64, error) {
	var dbError error
	lessonIds := make([]int64, 0)
	transaction, err := lessonsRepo.repository.startTransaction()
	if err != nil {
		return lessonIds, err
	}
	rows, err := transaction.Query(`SELECT id FROM lessons WHERE master_id = $1 AND date = $2 AND ((time_end > $3 AND time_start < $3) OR (time_start < $4 AND time_end > $4) OR (time_start = $3 AND time_end = $4)) `,
		lesson.MasterId, lesson.Date, lesson.TimeStart, lesson.TimeEnd)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve lessons: %v", err.Error())
		logger.Errorf(dbError.Error())
		return lessonIds, dbError
	}
	for rows.Next() {
		var lessFoundId int64
		err = rows.Scan(&lessFoundId)
		if err != nil {
			dbError = fmt.Errorf("failed to retrieve lesson: %v", err)
			logger.Errorf(dbError.Error())
			return lessonIds, dbError
		}
		lessonIds = append(lessonIds, lessFoundId)
	}
	err = lessonsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return lessonIds, err
	}
	return lessonIds, nil
}

func (lessonsRepo *LessonsRepo) GetLessonByIdAndMasterId(lesson *models.LessonDB, lessonId int64, masterId int64) error {
	var dbError error
	transaction, err := lessonsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT * FROM lessons WHERE id=$1 AND master_id=$2", lessonId, masterId)
	err = row.Scan(&lesson.Id, &lesson.MasterId, &lesson.TimeStart,
		&lesson.TimeEnd, &lesson.Date, &lesson.Price, &lesson.EducationFormat, &lesson.Status)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve lesson: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = lessonsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}
func (lessonsRepo *LessonsRepo) UpdateLessonByIdAndMasterId(lesson *models.LessonDB) error {
	var dbError error
	transaction, err := lessonsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	_, err = transaction.Exec("UPDATE lessons SET (time_start, time_end, date, price, education_format, status) = (coalesce($3, time_start), coalesce($4, time_end), coalesce($5, date), coalesce($6, price), coalesce($7, education_format), coalesce($8, status)) WHERE id = $1 AND master_id = $2",
		lesson.Id, lesson.MasterId, lesson.TimeStart, lesson.TimeEnd, lesson.Date, lesson.Price, lesson.EducationFormat, lesson.Status)
	if err != nil {
		dbError = fmt.Errorf("failed to update lesson: %v", err.Error())
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

//UPDATE lessons SET (time_start, time_end, date, price, education_format, status) = ('02:01:01', '02:01:02', '12-12-12', 100, 2, 1) WHERE id = 2 AND master_id = 1
