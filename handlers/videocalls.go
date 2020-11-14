package handlers

import "github.com/technoZoomers/MasterHubBackend/useCases"

type VCHandlers struct {
	handlers     *Handlers
	videocallsUC useCases.VideocallsUCInterface
	upgrader     websocket.Upgrader
}
