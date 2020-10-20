package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type LanguagesRepoI interface {
	GetAllLanguages() ([]string, error)
	GetLanguageById(language *models.LanguageDB) error
	GetLanguageByName(language *models.LanguageDB) error
}
