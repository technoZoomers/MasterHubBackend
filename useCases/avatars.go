package useCases

import "github.com/technoZoomers/MasterHubBackend/repository"

type AvatarsUC struct {
	useCases    *UseCases
	AvatarsRepo repository.AvatarsRepoI
}
