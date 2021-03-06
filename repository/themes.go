package repository

import (
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"strings"
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

func (themesRepo *ThemesRepo) searchIds(query string, source string, themeIds []int64, queryType string) ([]int64, error) {
	var dbError error
	themes := make([]int64, 0)
	transaction, err := themesRepo.repository.startTransaction()
	if err != nil {
		return themes, err
	}
	var queryValues []interface{}
	queryValues = append(queryValues, strings.ToLower(query))
	queryString := fmt.Sprintf(`SELECT id FROM %s WHERE name LIKE '%%' || $1 || '%%'`, source)
	if source == "subthemes" && len(themeIds) > 0 {
		queryString += fmt.Sprintf(" %s theme_id in (", queryType)
		queryCount := 1
		for _, th := range themeIds {
			queryCount++
			queryString += fmt.Sprintf("$%d,", queryCount)
			queryValues = append(queryValues, th)
		}
		queryString = queryString[:len(queryString)-1]
		queryString += ")"
	}
	rows, err := transaction.Query(queryString, queryValues...)
	if err != nil {
		dbError = fmt.Errorf("failed to retrieve %s: %v", source, err.Error())
		logger.Errorf(dbError.Error())
		return themes, dbError
	}
	for rows.Next() {
		var themeFoundId int64
		err = rows.Scan(&themeFoundId)
		if err != nil {
			dbError = fmt.Errorf("failed to retrieve one of the %s: %v", source, err)
			logger.Errorf(dbError.Error())
			return themes, dbError
		}
		themes = append(themes, themeFoundId)
	}
	err = themesRepo.repository.commitTransaction(transaction)
	if err != nil {
		return themes, err
	}
	return themes, nil
}

func (themesRepo *ThemesRepo) SearchThemeIds(query string) ([]int64, error) {
	return themesRepo.searchIds(query, "themes", []int64{}, "")
}

func (themesRepo *ThemesRepo) SearchSubthemeIds(query string) ([]int64, error) {
	return themesRepo.searchIds(query, "subthemes", []int64{}, "")
}

func (themesRepo *ThemesRepo) SearchSubthemeIdsAndThemes(query string, themeIds []int64) ([]int64, error) {
	return themesRepo.searchIds(query, "subthemes", themeIds, "AND")
}

func (themesRepo *ThemesRepo) SearchSubthemeIdsOrThemes(query string, themeIds []int64) ([]int64, error) {
	return themesRepo.searchIds(query, "subthemes", themeIds, "OR")
}