package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type LanguagesRepoI interface {
	GetAllLanguages() ([]string, error)
	GetLanguageById(language *models.LanguageDB) (int64, error)
}
