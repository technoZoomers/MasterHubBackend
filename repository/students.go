package repository

import (
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
)

type StudentsRepo struct {
	repository *Repository
}
func (studentsRepo *StudentsRepo) GetStudentByUserId(student *models.StudentDB) error {
	var dbError error
	transaction, err := studentsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT * FROM students WHERE user_id=$1", student.UserId)
	err = row.Scan(&student.Id, &student.UserId, &student.Username, &student.Fullname)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve student: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = studentsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}