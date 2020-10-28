package handlers

import (
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/gorilla/websocket"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"net/http"
)

type WSHandlers struct {
	handlers *Handlers
	ChatsUC  useCases.ChatsUCInterface
	upgrader websocket.Upgrader
	WebsocketsUC useCases.WebsocketsUCInterface
}

func (wsHandlers *WSHandlers) UpgradeConnection (writer http.ResponseWriter, req *http.Request) {
	sent, userId := wsHandlers.handlers.validateUserId(writer, req)
	if sent {
		return
	}
	//sent = wsHandlers.handlers.checkUserAuth(writer, req, userId)
	//if sent {
	//	return
	//}
	connection, err := wsHandlers.upgrader.Upgrade(writer, req, nil) // TODO: response header add
	if err != nil {
		upgradeError := fmt.Errorf("error upgrading connection: %v", err.Error())
		logger.Errorf(upgradeError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(upgradeError.Error()))
		return
	}
	websocketConnection := &models.WebsocketConnection{
		UserId:userId,
		Connection: connection,
	}
	wsHandlers.listenClient(websocketConnection)
}

func (wsHandlers *WSHandlers) listenClient(clientConn *models.WebsocketConnection) {
	wsHandlers.WebsocketsUC.AddClient(clientConn)
	for {
		_, message, err := clientConn.Connection.ReadMessage()
		if err != nil {
			readError := fmt.Errorf("error reading message from ws: %v", err.Error())
			logger.Errorf(readError.Error())
			wsHandlers.WebsocketsUC.RemoveClient(clientConn)
		}
		var wsMessage models.WebsocketMessage
		err = wsMessage.UnmarshalJSON(message)
		if err != nil {
			wsHandlers.WebsocketsUC.RemoveClient(clientConn)
		}
		if clientConn.UserId == wsMessage.Message.AuthorId {
			wsHandlers.WebsocketsUC.SendMessage(wsMessage)
		}
	}
}