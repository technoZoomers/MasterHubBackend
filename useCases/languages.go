package useCases

import (
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
)

type LanguagesUC struct {
	LanguagesRepo repository.LanguagesRepoI
}

func (languagesUC *LanguagesUC) Get() (models.Languages, error) {
	languages := make(models.Languages, 0)
	languages, err := languagesUC.LanguagesRepo.GetAllLanguages()
	if err != nil {
		return languages, err
	}
	return languages, nil
}
