package repository

import (
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/utils"
)

type ThemesRepo struct {
}

func (themesRepo *ThemesRepo) GetAllThemes() ([]models.ThemeDB, error) {
	themes := make([]models.ThemeDB, 0)
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return themes, dbError
	}
	rows, err := transaction.Query(`SELECT * FROM themes`)
	if err != nil {
		dbError := fmt.Errorf("failed to retrieve transactions: %v", err.Error())
		logger.Errorf(dbError.Error())
		return themes, dbError
	}
	for rows.Next() {
		var themeFound models.ThemeDB
		err = rows.Scan(&themeFound.Id, &themeFound.Name)
		if err != nil {
			logger.Errorf("failed to retrieve transaction: %v", err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				logger.Errorf("failed to rollback: %v", err)
				return themes, errRollback
			}
			return themes, err
		}
		themes = append(themes, themeFound)
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return themes, err
	}
	return themes, nil
}

func (themesRepo *ThemesRepo) GetSubthemesByTheme(theme *models.ThemeDB) ([]string, error) {
	subthemes := make([]string, 0)
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return subthemes, dbError
	}
	rows, err := transaction.Query(`SELECT name FROM subthemes WHERE theme_id=$1`, theme.Id)
	if err != nil {
		return subthemes, nil
	}
	for rows.Next() {
		var subthemeFound string
		err = rows.Scan(&subthemeFound)
		if err != nil {
			logger.Errorf("failed to retrieve transaction: %v", err)
			errRollback := transaction.Rollback()
			if errRollback != nil {
				logger.Errorf("failed to rollback: %v", err)
				return subthemes, errRollback
			}
			return subthemes, err
		}
		subthemes = append(subthemes, subthemeFound)
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return subthemes, err
	}
	return subthemes, nil
}

func (themesRepo *ThemesRepo) GetThemeById(theme *models.ThemeDB) (int64, error) {
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return utils.SERVER_ERROR, dbError
	}
	row := transaction.QueryRow("SELECT name FROM themes WHERE id=$1", theme.Id)
	err = row.Scan(&theme.Name)
	if err != nil {
		logger.Errorf("failed to retrieve theme: %v", err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			logger.Errorf("failed to rollback: %v", err)
			return utils.SERVER_ERROR, errRollback
		}
		return utils.USER_ERROR, fmt.Errorf("this theme doesn't exist")
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return utils.SERVER_ERROR, dbError
	}
	return utils.NO_ERROR, nil
}

func (themesRepo *ThemesRepo) GetSubthemeById (subtheme *models.SubthemeDB) (int64, error) {
	db := getPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return utils.SERVER_ERROR, dbError
	}
	row := transaction.QueryRow("SELECT name FROM subthemes WHERE id=$1", subtheme.Id)
	err = row.Scan(&subtheme.Name)
	if err != nil {
		logger.Errorf("failed to retrieve theme: %v", err)
		errRollback := transaction.Rollback()
		if errRollback != nil {
			logger.Errorf("failed to rollback: %v", err)
			return utils.SERVER_ERROR, errRollback
		}
		return utils.USER_ERROR, fmt.Errorf("this subtheme doesn't exist")
	}
	err = transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return utils.SERVER_ERROR, dbError
	}
	return utils.NO_ERROR, nil
}