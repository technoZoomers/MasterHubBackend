package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type ThemesRepoI interface {
	GetAllThemes() ([]models.ThemeDB, error)
	GetSubthemesByTheme(theme *models.ThemeDB) ([]string, error)
	GetThemeById(theme *models.ThemeDB) (int64, error)
	GetSubthemeById(subtheme *models.SubthemeDB) (int64, error)
}

