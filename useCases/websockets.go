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

func (wsUC *WebsocketsUC) Start() {
	for {
		select {
		case conn := <-wsUC.WebsocketsRepo.NewClients:
			wsUC.WebsocketsRepo.AddNewCLient(conn)
		case conn := <-wsUC.WebsocketsRepo.DroppedClients:
			wsUC.WebsocketsRepo.RemoveClient(conn)
		case message := <-wsUC.WebsocketsRepo.Messages:
			wsUC.processMessage(message)
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
	err := client.Connection.WriteMessage(websocket.TextMessage, marshalledMessage)
	if err != nil {
		broadcastError := fmt.Errorf("error broadcasting message: %v", err.Error())
		logger.Errorf(broadcastError.Error())
		return broadcastError
	}
	return nil
}
