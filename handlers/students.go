package handlers

import "github.com/technoZoomers/MasterHubBackend/useCases"

type StudentsHandlers struct {
	handlers   *Handlers
	StudentsUC useCases.StudentsUCInterface
}
