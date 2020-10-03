package repository

import (
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/utils"
)

type LanguagesRepo struct {
}

func (languagesRepo *LanguagesRepo) GetAllLanguages() ([]string, error) {
	languages := make([]string, 0)
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("Failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return languages, dbError
	}
	rows, err := transaction.Query(`SELECT name FROM languages`)
	if err != nil {
		dbError := fmt.Errorf("Failed to retrieve transactions: %v", err.Error())
		logger.Errorf(dbError.Error())
		return languages, dbError
	}
	for rows.Next() {
		var langFound string
		err = rows.Scan(&langFound)
		if err != nil {
			logger.Errorf("Failed to retrieve transaction: %v", err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				logger.Errorf("Failed to rollback: %v", err)
				return languages, errRollback
			}
			return languages, err
		}
		languages = append(languages, langFound)
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("Error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return languages, err
	}
	return languages, nil
}

func (languagesRepo *LanguagesRepo) GetLanguageById (language *models.LanguageDB) (int64, error) {
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("Failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return utils.SERVER_ERROR, dbError
	}
	row := transaction.QueryRow("SELECT name FROM languages WHERE id=$1", language.Id)
	err = row.Scan(&language.Name)
	if err != nil {
		logger.Errorf("Failed to retrieve theme: %v", err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			logger.Errorf("Failed to rollback: %v", err)
			return utils.SERVER_ERROR, errRollback
		}
		return utils.USER_ERROR, fmt.Errorf("this language doesn't exist")
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("Error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return utils.SERVER_ERROR, dbError
	}
	return utils.NO_ERROR, nil
}