package repository

import (
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
)

type ThemesRepo struct {
	repository *Repository
}

func (themesRepo *ThemesRepo) GetAllThemes() ([]models.ThemeDB, error) {
	var dbError error
	themes := make([]models.ThemeDB, 0)
	transaction, err := themesRepo.repository.startTransaction()
	if err != nil {
		return themes, err
	}
	rows, err := transaction.Query(`SELECT * FROM themes`)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve themes: %v", err.Error())
		logger.Errorf(dbError.Error())
		return themes, dbError
	}
	for rows.Next() {
		var themeFound models.ThemeDB
		err = rows.Scan(&themeFound.Id, &themeFound.Name)
		if err != nil {
			dbError = fmt.Errorf("failed to retrieve theme: %v", err)
			logger.Errorf(dbError.Error())
			return themes, dbError
		}
		themes = append(themes, themeFound)
	}
	err = themesRepo.repository.commitTransaction(transaction)
	if err != nil {
		return themes, err
	}
	return themes, nil
}

func (themesRepo *ThemesRepo) GetSubthemesByTheme(theme *models.ThemeDB) ([]string, error) {
	var dbError error
	subthemes := make([]string, 0)
	transaction, err := themesRepo.repository.startTransaction()
	if err != nil {
		return subthemes, err
	}
	rows, err := transaction.Query(`SELECT name FROM subthemes WHERE theme_id=$1`, theme.Id)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve subthemes by theme: %v", err.Error())
		logger.Errorf(dbError.Error())
		return subthemes, dbError
	}
	for rows.Next() {
		var subthemeFound string
		err = rows.Scan(&subthemeFound)
		if err != nil {
			dbError = fmt.Errorf("failed to retrieve subtheme: %v", err)
			logger.Errorf(dbError.Error())
			return subthemes, dbError
		}
		subthemes = append(subthemes, subthemeFound)
	}
	err = themesRepo.repository.commitTransaction(transaction)
	if err != nil {
		return subthemes, err
	}
	return subthemes, nil
}

func (themesRepo *ThemesRepo) GetThemeById(theme *models.ThemeDB) error {
	var dbError error
	transaction, err := themesRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT name FROM themes WHERE id=$1", theme.Id)
	err = row.Scan(&theme.Name)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve theme: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = themesRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (themesRepo *ThemesRepo) GetThemeByName(theme *models.ThemeDB) error {
	var dbError error
	transaction, err := themesRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT id FROM themes WHERE name=$1", theme.Name)
	err = row.Scan(&theme.Id)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve theme: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = themesRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (themesRepo *ThemesRepo) GetSubthemeById(subtheme *models.SubthemeDB) error {
	var dbError error
	transaction, err := themesRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT name FROM subthemes WHERE id=$1", subtheme.Id)
	err = row.Scan(&subtheme.Name)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve subtheme: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = themesRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (themesRepo *ThemesRepo) GetSubthemeByName(subtheme *models.SubthemeDB) error {
	var dbError error
	transaction, err := themesRepo.repository.startTransaction()
	if err != nil {
		return err
	}
	row := transaction.QueryRow("SELECT id, theme_id FROM subthemes WHERE name=$1", subtheme.Name)
	err = row.Scan(&subtheme.Id, &subtheme.ThemeId)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve subtheme: %v", err.Error())
		logger.Errorf(dbError.Error())
	}
	err = themesRepo.repository.commitTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}
