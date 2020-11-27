package useCases

import "github.com/technoZoomers/MasterHubBackend/models"

type WebsocketsUCInterface interface {
	AddClient(clientConnection *models.WebsocketConnection)
	RemoveClient(clientConnection *models.WebsocketConnection)
	SendMessage(message models.WebsocketMessage)
	SendNotification(notification models.WebsocketNotification)
	CheckOnline(userId int64) bool
	Start()
}
