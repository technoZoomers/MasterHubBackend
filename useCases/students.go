package useCases

import "github.com/technoZoomers/MasterHubBackend/repository"

type StudentsUC struct {
	useCases     *UseCases
	StudentsRepo repository.StudentsRepoI
}
