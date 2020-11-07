package repository

import (
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
)

type LanguagesRepo struct {
	repository *Repository
}

func (languagesRepo *LanguagesRepo) GetAllLanguages() ([]string, error) {
	var dbError error
	languages := make([]string, 0)
	transaction, err := languagesRepo.repository.startTransaction()
	if err != nil {
		return languages, err
	}
	rows, err := transaction.Query(`SELECT name FROM languages`)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve languages: %v", err.Error())
		logger.Errorf(dbError.Error())
		return languages, dbError
	}
	for rows.Next() {
		var langFound string
		err = rows.Scan(&langFound)
		if err != nil {
			dbError = fmt.Errorf("failed to retrieve language: %v", err)
			logger.Errorf(dbError.Error())
			return languages, dbError
		}
		languages = append(languages, langFound)
	}
	err = languagesRepo.repository.commitTransaction(transaction)
	if err != nil {
		return languages, err
	}
	return languages, nil
}

func (languagesRepo *LanguagesRepo) GetLanguageById(language *models.LanguageDB) error {
	var dbError error
	transaction, err := languagesRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT name FROM languages WHERE id=$1", language.Id)
	err = row.Scan(&language.Name)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve language: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = languagesRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (languagesRepo *LanguagesRepo) GetLanguageByName(language *models.LanguageDB) error {
	var dbError error
	transaction, err := languagesRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT id FROM languages WHERE name=$1", language.Name)
	err = row.Scan(&language.Id)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve language: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = languagesRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}
