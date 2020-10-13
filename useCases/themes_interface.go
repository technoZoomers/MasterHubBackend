package useCases

import "github.com/technoZoomers/MasterHubBackend/models"

type ThemesUCInterface interface {
	Get() (models.Themes, error)
	GetThemeById(theme *models.Theme) (bool, error)
}
