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

func (lessonsRepo *LessonsRepo) InsertLessonRequest(lessonRequest *models.LessonStudentDB) error {
	var dbError error
	transaction, err := lessonsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	_, err = transaction.Exec("INSERT INTO lessons_students (lesson_id, student_id, status) VALUES ($1, $2, $3)", lessonRequest.LessonId, lessonRequest.StudentId, lessonRequest.Status)
	if err != nil {
		dbError = fmt.Errorf("failed to insert lesson request: %v", err.Error())
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
func (lessonsRepo *LessonsRepo) GetLessonRequestByStudentIdAndLessonId(lessonRequest *models.LessonStudentDB, studentId int64, lessonId int64) error {
	var dbError error
	transaction, err := lessonsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT * FROM lessons_students WHERE lesson_id=$1 AND student_id=$2", lessonId, studentId)
	err = row.Scan(&lessonRequest.LessonId, &lessonRequest.StudentId, &lessonRequest.Status)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve lesson student: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = lessonsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (lessonsRepo *LessonsRepo) GetLessonRequestByStudentUserIdAndLessonId(lessonRequest *models.LessonStudentDB, studentId int64, lessonId int64) error {
	var dbError error
	transaction, err := lessonsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT * FROM lessons_students WHERE lesson_id=$1 AND student_id in (select id from students where user_id=$2)", lessonId, studentId)
	err = row.Scan(&lessonRequest.LessonId, &lessonRequest.StudentId, &lessonRequest.Status)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve lesson student: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = lessonsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (lessonsRepo *LessonsRepo) DeleteLessonRequestByStudentIdAndLessonId(studentId int64, lessonId int64) error {
	var dbError error
	transaction, err := lessonsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	_, err = transaction.Exec(" DELETE FROM lessons_students WHERE lesson_id=$1 AND student_id=$2", lessonId, studentId)
	if err != nil {
		dbError = fmt.Errorf("failed to delete request: %v", err.Error())
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

func (lessonsRepo *LessonsRepo) GetMastersLessons(masterId int64, query models.LessonsQueryValuesDB) ([]models.LessonDB, error) {
	var dbError error
	lessons := make([]models.LessonDB, 0)
	transaction, err := lessonsRepo.repository.startTransaction()
	if err != nil {
		return lessons, err
	}
	var selectQueryString string
	var queryValues []interface{}
	queryValues = append(queryValues, masterId)
	if query.Status != 0 {
		selectQueryString = `SELECT * FROM lessons WHERE master_id = $1 and status = $2`
		queryValues = append(queryValues, query.Status)
	} else {
		selectQueryString = `SELECT * FROM lessons WHERE master_id = $1`
	}
	rows, err := transaction.Query(selectQueryString, queryValues...)
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

func (lessonsRepo *LessonsRepo) GetStudentsLessons(studentId int64, query models.LessonsQueryValuesDB) ([]models.LessonDB, error) {
	var dbError error
	lessons := make([]models.LessonDB, 0)
	transaction, err := lessonsRepo.repository.startTransaction()
	if err != nil {
		return lessons, err
	}
	var selectQueryString string
	var queryValues []interface{}
	queryValues = append(queryValues, studentId)
	if query.Status != 0 {
		selectQueryString = `select lessons.id, time_start, time_end, date, price, lessons.education_format, lessons.status, masters.user_id from  lessons join lessons_students on lessons.id = lessons_students.lesson_id join masters on lessons.master_id = masters.id where student_id = $1 and lessons.status = $2`
		queryValues = append(queryValues, query.Status)
	} else {
		selectQueryString = `select lessons.id, time_start, time_end, date, price, lessons.education_format, lessons.status, masters.user_id from  lessons join lessons_students on lessons.id = lessons_students.lesson_id join masters on lessons.master_id = masters.id where student_id = $1`
	}
	rows, err := transaction.Query(selectQueryString, queryValues...)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve lessons: %v", err.Error())
		logger.Errorf(dbError.Error())
		return lessons, dbError
	}
	for rows.Next() {
		var lessFound models.LessonDB
		err = rows.Scan(&lessFound.Id, &lessFound.TimeStart,
			&lessFound.TimeEnd, &lessFound.Date, &lessFound.Price, &lessFound.EducationFormat, &lessFound.Status, &lessFound.MasterId)
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
func (lessonsRepo *LessonsRepo) GetLessonById(lesson *models.LessonDB, lessonId int64) error {
	var dbError error
	transaction, err := lessonsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT * FROM lessons WHERE id=$1", lessonId)
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

func (lessonsRepo *LessonsRepo) DeleteLessonById(lessonId int64) error {
	var dbError error
	transaction, err := lessonsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	_, err = transaction.Exec(" DELETE FROM lessons WHERE id=$1", lessonId)
	if err != nil {
		dbError = fmt.Errorf("failed to delete lesson: %v", err.Error())
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

func (lessonsRepo *LessonsRepo) GetLessonStudents(lessonId int64) ([]int64, error) {
	var dbError error
	lessonStudents := make([]int64, 0)
	transaction, err := lessonsRepo.repository.startTransaction()
	if err != nil {
		return lessonStudents, err
	}
	rows, err := transaction.Query(`select user_id from lessons_students join students on lessons_students.student_id = students.id where lesson_id=$1 AND status=1`, lessonId)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve lesson students: %v", err.Error())
		logger.Errorf(dbError.Error())
		return lessonStudents, dbError
	}
	for rows.Next() {
		var studentId int64
		err = rows.Scan(&studentId)
		if err != nil {
			dbError = fmt.Errorf("failed to retrieve student: %v", err)
			logger.Errorf(dbError.Error())
			return lessonStudents, dbError
		}
		lessonStudents = append(lessonStudents, studentId)
	}
	err = lessonsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return lessonStudents, err
	}
	return lessonStudents, nil
}

func (lessonsRepo *LessonsRepo) GetMastersLessonsRequests(masterId int64, query models.LessonsQueryValuesDB) ([]models.LessonStudentDB, error) {
	var dbError error
	lessonStudents := make([]models.LessonStudentDB, 0)
	transaction, err := lessonsRepo.repository.startTransaction()
	if err != nil {
		return lessonStudents, err
	}
	var selectQueryString string
	var queryValues []interface{}
	queryValues = append(queryValues, masterId)
	if query.Status != 0 {
		selectQueryString = `select lesson_id, user_id, lessons_students.status from lessons_students join students on lessons_students.student_id = students.id join lessons on lessons_students.lesson_id = lessons.id where master_id=$1 and lessons_students.status=$2`
		queryValues = append(queryValues, query.Status)
	} else {
		selectQueryString = `select lesson_id, user_id, lessons_students.status from lessons_students join students on lessons_students.student_id = students.id join lessons on lessons_students.lesson_id = lessons.id where master_id=$1`
	}
	rows, err := transaction.Query(selectQueryString, queryValues...)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve master lesson requests: %v", err.Error())
		logger.Errorf(dbError.Error())
		return lessonStudents, dbError
	}
	for rows.Next() {
		var lessonRequest models.LessonStudentDB
		err = rows.Scan(&lessonRequest.LessonId, &lessonRequest.StudentId, &lessonRequest.Status)
		if err != nil {
			dbError = fmt.Errorf("failed to retrieve request: %v", err)
			logger.Errorf(dbError.Error())
			return lessonStudents, dbError
		}
		lessonStudents = append(lessonStudents, lessonRequest)
	}
	err = lessonsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return lessonStudents, err
	}
	return lessonStudents, nil
}

func (lessonsRepo *LessonsRepo) GetStudentsLessonsRequests(studentId int64, query models.LessonsQueryValuesDB) ([]models.LessonStudentDB, error) {
	var dbError error
	lessonStudents := make([]models.LessonStudentDB, 0)
	transaction, err := lessonsRepo.repository.startTransaction()
	if err != nil {
		return lessonStudents, err
	}
	var selectQueryString string
	var queryValues []interface{}
	queryValues = append(queryValues, studentId)
	if query.Status != 0 {
		selectQueryString = `select lesson_id, student_id, status from lessons_students where student_id=$1 and status=$2`
		queryValues = append(queryValues, query.Status)
	} else {
		selectQueryString = `select lesson_id, student_id, status from lessons_students where student_id=$1`
	}
	rows, err := transaction.Query(selectQueryString, queryValues...)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve student lesson requests: %v", err.Error())
		logger.Errorf(dbError.Error())
		return lessonStudents, dbError
	}
	for rows.Next() {
		var lessonRequest models.LessonStudentDB
		err = rows.Scan(&lessonRequest.LessonId, &lessonRequest.StudentId, &lessonRequest.Status)
		if err != nil {
			dbError = fmt.Errorf("failed to retrieve request: %v", err)
			logger.Errorf(dbError.Error())
			return lessonStudents, dbError
		}
		lessonStudents = append(lessonStudents, lessonRequest)
	}
	err = lessonsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return lessonStudents, err
	}
	return lessonStudents, nil
}

func (lessonsRepo *LessonsRepo) UpdateLessonRequest(lessonRequest *models.LessonStudentDB) error {
	var dbError error
	transaction, err := lessonsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	_, err = transaction.Exec("update lessons_students set status=$1 where lesson_id=$2 and student_id=$3",
		lessonRequest.Status, lessonRequest.LessonId, lessonRequest.StudentId)
	if err != nil {
		dbError = fmt.Errorf("failed to update lesson request: %v", err.Error())
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
