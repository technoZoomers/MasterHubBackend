package useCases

import "github.com/technoZoomers/MasterHubBackend/models"

type WebsocketsUCInterface interface {
	AddClient(clientConnection *models.WebsocketConnection)
	RemoveClient(clientConnection *models.WebsocketConnection)
	SendMessage(message models.WebsocketMessage)
	Start()
}
