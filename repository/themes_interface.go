package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type ThemesRepoI interface {
	GetAllThemes() ([]models.ThemeDB, error)
	GetSubthemesByTheme(theme *models.ThemeDB) ([]string, error)
	GetThemeById(theme *models.ThemeDB) error
	GetThemeByName(theme *models.ThemeDB) error
	GetSubthemeById(subtheme *models.SubthemeDB) error
	GetSubthemeByName(subtheme *models.SubthemeDB) error
	SearchSubthemeIds(query string) ([]int64, error)
	SearchThemeIds(query string) ([]int64, error)
	SearchSubthemeIdsAndThemes(query string, themeIds []int64) ([]int64, error)
	SearchSubthemeIdsOrThemes(query string, themeIds []int64) ([]int64, error)
}
