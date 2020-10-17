package useCases

import (
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
)

type ThemesUC struct {
	useCases   *UseCases
	ThemesRepo repository.ThemesRepoI
}

func (themesUC *ThemesUC) Get() (models.Themes, error) {
	themes := make(models.Themes, 0)
	themesDB, err := themesUC.ThemesRepo.GetAllThemes()
	if err != nil {
		return themes, fmt.Errorf(themesUC.useCases.errorMessages.DbError)
	}
	for _, theme := range themesDB {
		subthemesDB, err := themesUC.ThemesRepo.GetSubthemesByTheme(&theme)
		if err != nil {
			return themes, fmt.Errorf(themesUC.useCases.errorMessages.DbError)
		}
		themes = append(themes, models.Theme{
			Id:        theme.Id,
			Theme:     theme.Name,
			Subthemes: subthemesDB,
		})
	}
	return themes, nil
}

func (themesUC *ThemesUC) GetThemeById(theme *models.Theme) error {
	if theme.Id == themesUC.useCases.errorId {
		return &models.BadRequestError{Message: "incorrect theme id", RequestId: theme.Id}
	}
	var themeDB models.ThemeDB
	themeDB.Id = theme.Id
	err := themesUC.ThemesRepo.GetThemeById(&themeDB)
	if err != nil {
		return fmt.Errorf(themesUC.useCases.errorMessages.DbError)
	}
	if themeDB.Name == "" {
		absenceError := &models.BadRequestError{Message: "theme doesn't exist", RequestId: theme.Id}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	theme.Theme = themeDB.Name
	subthemesDB, err := themesUC.ThemesRepo.GetSubthemesByTheme(&themeDB)
	if err != nil {
		return fmt.Errorf(themesUC.useCases.errorMessages.DbError)
	}
	theme.Subthemes = subthemesDB
	return nil
}
