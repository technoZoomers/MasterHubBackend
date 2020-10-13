package useCases

import "github.com/technoZoomers/MasterHubBackend/models"

type LanguagesUCInterface interface {
	Get() (models.Languages, error)
}
