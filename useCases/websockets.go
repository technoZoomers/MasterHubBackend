package useCases

import (
	"fmt"
	"github.com/google/logger"
	"github.com/gorilla/websocket"
	json "github.com/mailru/easyjson"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
)

type WebsocketsUC struct {
	useCases         *UseCases
	WebsocketsRepo   repository.WebsocketsRepo // TODO: INTERFACE
	ChatsRepo        repository.ChatsRepoI
	messagesTypesMap map[int64]bool
}

func (wsUC *WebsocketsUC) AddClient(clientConnection *models.WebsocketConnection) {
	wsUC.WebsocketsRepo.NewClients <- clientConnection
}

func (wsUC *WebsocketsUC) RemoveClient(clientConnection *models.WebsocketConnection) {
	wsUC.WebsocketsRepo.DroppedClients <- clientConnection
}

func (wsUC *WebsocketsUC) SendMessage(message models.WebsocketMessage) {
	wsUC.WebsocketsRepo.Messages <- message
}

func (wsUC *WebsocketsUC) SendNotification(notification models.WebsocketNotification) {
	wsUC.WebsocketsRepo.Notifications <- notification
}

func (wsUC *WebsocketsUC) CheckOnline(userId int64) bool {
	userConnString := wsUC.WebsocketsRepo.GetConnectionString(userId)
	if userConnString != "" {
		return false
	}
	return true
}

func (wsUC *WebsocketsUC) Start() {
	for {
		select {
		case conn := <-wsUC.WebsocketsRepo.NewClients:
			wsUC.WebsocketsRepo.AddNewCLient(conn)
		case conn := <-wsUC.WebsocketsRepo.DroppedClients:
			wsUC.WebsocketsRepo.RemoveClient(conn)
		case message := <-wsUC.WebsocketsRepo.Messages:
			wsUC.processMessage(message)
		case notification := <-wsUC.WebsocketsRepo.Notifications:
			wsUC.processNotification(notification)
		}
	}
}

func (wsUC *WebsocketsUC) matchMessageToDB(messageDB *models.MessageDB, message *models.Message) {
	messageDB.UserId = message.AuthorId
	messageDB.ChatId = message.ChatId
	messageDB.Created = message.Created
	messageDB.Text = message.Text
	messageDB.Info = wsUC.messagesTypesMap[message.Type]
}

func (wsUC *WebsocketsUC) processNotification(notification models.WebsocketNotification) {
	marshalledNotification, err := json.Marshal(notification)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Errorf(jsonError.Error())
		return
	}
	if notification.Notification.UserId != 0 {
		userConnString := wsUC.WebsocketsRepo.GetConnectionString(notification.Notification.UserId)
		if userConnString != "" {
			err = wsUC.writeMessageToConnection(userConnString, marshalledNotification)
			if err != nil {
				return
			}
		}
	}
}

func (wsUC *WebsocketsUC) processMessage(message models.WebsocketMessage) {
	marshalledMessage, err := json.Marshal(message)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Errorf(jsonError.Error())
		return
	}

	var chat models.ChatDB
	err = wsUC.ChatsRepo.GetChatById(&chat, message.Message.ChatId)
	if err != nil {
		dbError := fmt.Errorf(wsUC.useCases.errorMessages.DbError)
		logger.Errorf(dbError.Error())
		return
	}

	if message.Type == 1 { // TODO: other message types 2-delete, 3-resend
		var messageDB models.MessageDB
		wsUC.matchMessageToDB(&messageDB, &message.Message)
		err = wsUC.ChatsRepo.InsertMessage(&messageDB)
		if err != nil {
			dbError := fmt.Errorf(wsUC.useCases.errorMessages.DbError)
			logger.Errorf(dbError.Error())
			return
		}
	}

	if chat.StudentId != 0 {
		studentConnString := wsUC.WebsocketsRepo.GetConnectionString(chat.StudentId)
		if studentConnString != "" {
			err = wsUC.writeMessageToConnection(studentConnString, marshalledMessage)
			if err != nil {
				return
			}
		}
	}
	if chat.MasterId != 0 {
		masterConnString := wsUC.WebsocketsRepo.GetConnectionString(chat.MasterId)
		if masterConnString != "" {
			err = wsUC.writeMessageToConnection(masterConnString, marshalledMessage)
			if err != nil {
				return
			}
		}
	}
}

func (wsUC *WebsocketsUC) writeMessageToConnection(connectionString string, marshalledMessage []byte) error {
	client := wsUC.WebsocketsRepo.GetConnection(connectionString)
	if client != nil {
		err := client.Connection.WriteMessage(websocket.TextMessage, marshalledMessage)
		if err != nil {
			broadcastError := fmt.Errorf("error broadcasting message: %v", err.Error())
			logger.Errorf(broadcastError.Error())
			return broadcastError
		}
	}
	return nil
}
