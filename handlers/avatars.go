package handlers

import "github.com/technoZoomers/MasterHubBackend/useCases"

type AvatarsHandlers struct {
	handlers  *Handlers
	AvatarsUC useCases.AvatarsUCInterface
}
