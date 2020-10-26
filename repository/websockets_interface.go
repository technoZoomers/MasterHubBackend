package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type WebsocketsRepoI interface {
	AddNewCLient(clientConnection *models.WebsocketConnection)
	RemoveClient(clientConnection *models.WebsocketConnection)
	GetConnection(connString string) *models.WebsocketConnection
	GetConnectionString(userId int64) string
}
