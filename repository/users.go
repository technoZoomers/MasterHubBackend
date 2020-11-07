package repository

import (
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
)

type UsersRepo struct {
	repository *Repository
}

func (usersRepo *UsersRepo) InsertUser(user *models.UserDB) error {
	var dbError error
	transaction, err := usersRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("INSERT INTO users (email, password, type, created) values ($1, $2, $3, $4) returning id",
		user.Email, user.Password, user.Type, user.Created)
	err = row.Scan(&user.Id)
	if err != nil {
		dbError = fmt.Errorf("failed to insert user: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := usersRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = usersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (usersRepo *UsersRepo) DeleteUserWithId(userId int64) error {
	var dbError error
	transaction, err := usersRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	_, err = transaction.Exec("DELETE FROM users WHERE id = $1", userId)
	if err != nil {
		dbError = fmt.Errorf("failed to delete user: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := usersRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = usersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (usersRepo *UsersRepo) GetUserById(user *models.UserDB, userId int64) error {
	var dbError error
	transaction, err := usersRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT * FROM users WHERE id=$1", userId)
	err = row.Scan(&user.Id, &user.Email, &user.Password, &user.Type, &user.Created)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve user: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = usersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (usersRepo *UsersRepo) GetUserByEmail(user *models.UserDB) error {
	var dbError error
	transaction, err := usersRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT * FROM users WHERE email=$1", user.Email)
	err = row.Scan(&user.Id, &user.Email, &user.Password, &user.Type, &user.Created)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve user: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = usersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (usersRepo *UsersRepo) GetUserByEmailAndPassword(user *models.UserDB) error {
	var dbError error
	transaction, err := usersRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT * FROM users WHERE email=$1 and password=$2", user.Email, user.Password)
	err = row.Scan(&user.Id, &user.Email, &user.Password, &user.Type, &user.Created)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve user: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = usersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (usersRepo *UsersRepo) GetUserLanguagesById(userId int64) ([]int64, error) {
	languagesIds := make([]int64, 0)
	var dbError error
	transaction, err := usersRepo.repository.startTransaction()
	if err != nil {
		return languagesIds, err
	}
	rows, err := transaction.Query(`SELECT language_id FROM users_languages WHERE user_id=$1`, userId)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve language ids: %v", err.Error())
		logger.Errorf(dbError.Error())
		return languagesIds, dbError
	}
	for rows.Next() {
		var languageIdFound int64
		err = rows.Scan(&languageIdFound)
		if err != nil {
			dbError = fmt.Errorf("failed to retrieve language id: %v", err)
			logger.Errorf(dbError.Error())
			return languagesIds, dbError
		}
		languagesIds = append(languagesIds, languageIdFound)
	}
	err = usersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return languagesIds, err
	}
	return languagesIds, nil
}

func (usersRepo *UsersRepo) DeleteUserLanguagesById(userId int64) error {
	var dbError error
	transaction, err := usersRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	_, err = transaction.Exec("DELETE FROM users_languages WHERE user_id=$1", userId)
	if err != nil {
		dbError = fmt.Errorf("failed to delete languages: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := usersRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = usersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (usersRepo *UsersRepo) SetUserLanguagesById(userId int64, languages []int64) error {
	if languages == nil || len(languages) == 0 {
		return nil
	}
	var dbError error
	transaction, err := usersRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	var queryValues []interface{}
	insertQuery := "INSERT INTO users_languages (user_id, language_id) values "
	queryValues = append(queryValues, userId)
	for i, subth := range languages {
		insertQuery += fmt.Sprintf("($1, $%d),", i+2)
		queryValues = append(queryValues, subth)
	}
	insertQuery = insertQuery[:len(insertQuery)-1]
	_, err = transaction.Exec(insertQuery, queryValues...)
	if err != nil {
		dbError = fmt.Errorf("failed to insert languages: %v", err.Error())
		logger.Errorf(dbError.Error())
		errRollback := usersRepo.repository.rollbackTransaction(transaction)
		if errRollback != nil {
			return errRollback
		}
		return dbError
	}
	err = usersRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}
