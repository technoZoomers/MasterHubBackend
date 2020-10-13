package useCases

import (
	"fmt"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
	"github.com/technoZoomers/MasterHubBackend/utils"
)

type ThemesUC struct {
	ThemesRepo repository.ThemesRepoI
}

func (themesUC *ThemesUC) Get() (models.Themes, error) {
	themes := make(models.Themes, 0)
	themesDB, err := themesUC.ThemesRepo.GetAllThemes()
	if err != nil {
		return themes, err
	}
	for _, theme := range themesDB {
		subthemesDB, err := themesUC.ThemesRepo.GetSubthemesByTheme(&theme)
		if err != nil {
			return themes, err
		}
		themes = append(themes, models.Theme{
			Id:        theme.Id,
			Theme:     theme.Name,
			Subthemes: subthemesDB,
		})
	}
	return themes, nil
}

func (themesUC *ThemesUC) GetThemeById(theme *models.Theme) (bool, error) {
	if theme.Id == utils.ERROR_ID {
		return true, fmt.Errorf("incorrect theme id")
	}
	var themeDB models.ThemeDB
	themeDB.Id = theme.Id
	errType, err := themesUC.ThemesRepo.GetThemeById(&themeDB)
	if err != nil {
		if errType == utils.USER_ERROR {
			return true, err
		} else if errType == utils.SERVER_ERROR {
			return false, err
		}
	}
	theme.Theme = themeDB.Name
	subthemesDB, err := themesUC.ThemesRepo.GetSubthemesByTheme(&themeDB)
	if err != nil {
		return false, err
	}
	theme.Subthemes = subthemesDB
	return false, nil
}
