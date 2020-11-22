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

func (studentsRepo *StudentsRepo) GetStudentByUserIdWithEmail(student *models.StudentDB) (string, error) {
	var dbError error
	var email string
	transaction, err := studentsRepo.repository.startTransaction()
	if err != nil {
		return email, err
	}
	row := transaction.QueryRow("select students.id, user_id, username, fullname, email from students join users u on students.user_id = u.id where students.user_id = $1", student.UserId)
	err = row.Scan(&student.Id, &student.UserId, &student.Username, &student.Fullname, &email)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve student: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = studentsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return email, err
	}
	return email, nil
}

func (studentsRepo *StudentsRepo) GetStudentIdByUsername(student *models.StudentDB) error {
	var dbError error
	transaction, err := studentsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT id FROM students WHERE username=$1", student.Username)
	err = row.Scan(&student.Id)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve student id: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = studentsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (studentsRepo *StudentsRepo) InsertStudent(student *models.StudentDB) error {
	var dbError error
	transaction, err := studentsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("INSERT INTO students (user_id, username, fullname) values ($1, $2, $3) returning id",
		student.UserId, student.Username, student.Fullname)
	err = row.Scan(&student.Id)
	if err != nil {
		dbError = fmt.Errorf("failed to insert student: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := studentsRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = studentsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (studentsRepo *StudentsRepo) UpdateStudent(student *models.StudentDB) error {
	var dbError error
	transaction, err := studentsRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	_, err = transaction.Exec("UPDATE students SET (username, fullname) = ($1, $2) where id = $3",
		student.Username, student.Fullname, student.Id)
	if err != nil {
		dbError = fmt.Errorf("failed to update student: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := studentsRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = studentsRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}
