package handlers

import "github.com/technoZoomers/MasterHubBackend/useCases"

type UsersHandlers struct {
	handlers     *Handlers
	UsersUC useCases.UsersUCInterface
}
