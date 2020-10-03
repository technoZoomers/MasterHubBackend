package useCases

import "github.com/technoZoomers/MasterHubBackend/repository"

type StudentsUC struct {
	StudentsRepo  repository.StudentsRepoI
}