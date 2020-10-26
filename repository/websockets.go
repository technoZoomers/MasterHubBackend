package repository

import (
	"github.com/technoZoomers/MasterHubBackend/models"
)

type WebsocketsRepo struct {
	repository *Repository
	userConnMap map[int64]string
	clientsMap          map[string]*models.WebsocketConnection
	NewClients    chan *models.WebsocketConnection
	DroppedClients chan *models.WebsocketConnection
	Messages    chan models.WebsocketMessage
}

func (wsRepo *WebsocketsRepo) AddNewCLient(clientConnection *models.WebsocketConnection) {
	connString := clientConnection.Connection.RemoteAddr().String()
	wsRepo.clientsMap[connString] = clientConnection
	wsRepo.userConnMap[clientConnection.UserId] = connString
}

func (wsRepo *WebsocketsRepo) RemoveClient(clientConnection *models.WebsocketConnection) {
	delete(wsRepo.clientsMap, clientConnection.Connection.RemoteAddr().String())
}

func (wsRepo *WebsocketsRepo)  GetConnection(connString string) *models.WebsocketConnection {
	return wsRepo.clientsMap[connString]
}

func (wsRepo *WebsocketsRepo)  GetConnectionString(userId int64) string {
	return wsRepo.userConnMap[userId]
}