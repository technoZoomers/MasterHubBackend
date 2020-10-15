package useCases

import "github.com/technoZoomers/MasterHubBackend/repository"

type UsersUC struct {
	useCases  *UseCases
	UsersRepo repository.UsersRepoI
}
