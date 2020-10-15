package useCases

import (
	"fmt"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
)

type LanguagesUC struct {
	useCases      *UseCases
	LanguagesRepo repository.LanguagesRepoI
}

func (languagesUC *LanguagesUC) Get() (models.Languages, error) {
	languages := make(models.Languages, 0)
	languages, err := languagesUC.LanguagesRepo.GetAllLanguages()
	if err != nil {
		return languages, fmt.Errorf(languagesUC.useCases.errorMessages.DbError)
	}
	return languages, nil
}
