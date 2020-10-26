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
	useCases     *UseCases
	WebsocketsRepo   repository.WebsocketsRepo// TODO: INTERFACE
	ChatsRepo   repository.ChatsRepoI
}


func (wsUC *WebsocketsUC) AddClient(clientConnection *models.WebsocketConnection) {
	wsUC.WebsocketsRepo.NewClients <-clientConnection
}

func (wsUC *WebsocketsUC) RemoveClient(clientConnection *models.WebsocketConnection) {
	wsUC.WebsocketsRepo.DroppedClients <-clientConnection
}

func (wsUC *WebsocketsUC) SendMessage(message models.WebsocketMessage) {
	wsUC.WebsocketsRepo.Messages <-message
}

func (wsUC *WebsocketsUC) Start() {
	for {
		select {
		case conn := <-wsUC.WebsocketsRepo.NewClients:
			wsUC.WebsocketsRepo.AddNewCLient(conn)
		case conn := <-wsUC.WebsocketsRepo.DroppedClients:
			wsUC.WebsocketsRepo.RemoveClient(conn)
		case message := <-wsUC.WebsocketsRepo.Messages:
			err := wsUC.processMessage(message)
			if err != nil {
				logger.Errorf(err.Error())
			}
		}
	}
}

func (wsUC *WebsocketsUC) matchMessageToDB (messageDB *models.MessageDB, message *models.Message) {
	messageDB.UserId = message.AuthorId
	messageDB.ChatId = message.ChatId
	messageDB.Created = message.Created
	messageDB.Text = message.Text
	if message.Type == 1 { //TODO: refactor types!!!
		messageDB.Info = false
	} else {
		messageDB.Info = true
	}
}

func (wsUC *WebsocketsUC) processMessage(message models.WebsocketMessage) error {
	marshalledMessage, err := json.Marshal(message)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Error(jsonError)
		return jsonError
	}

	var chat models.ChatDB
	err = wsUC.ChatsRepo.GetChatById(&chat, message.Message.ChatId)
	if err != nil {
		return fmt.Errorf(wsUC.useCases.errorMessages.DbError)
	}

	if message.Type == 1 {
		var messageDB models.MessageDB
		wsUC.matchMessageToDB(&messageDB, &message.Message)
		err = wsUC.ChatsRepo.InsertMessage(&messageDB)
		if err != nil {
			return fmt.Errorf(wsUC.useCases.errorMessages.DbError)
		}
	}

	if chat.StudentId != 0 {
		studentConnString := wsUC.WebsocketsRepo.GetConnectionString(chat.StudentId)
		if studentConnString != "" {
			client := wsUC.WebsocketsRepo.GetConnection(studentConnString)
			err := client.Connection.WriteMessage(websocket.TextMessage, marshalledMessage)
			if err != nil { // TODO: refactor!!!
				fmt.Println("Error broadcasting message: ", err)
				return err
			}
		}
	}
	if chat.MasterId != 0 {
		masterConnString := wsUC.WebsocketsRepo.GetConnectionString(chat.MasterId)
		if masterConnString != "" {
			client := wsUC.WebsocketsRepo.GetConnection(masterConnString)
			err := client.Connection.WriteMessage(websocket.TextMessage, marshalledMessage)
			if err != nil { // TODO: refactor!!!
				fmt.Println("Error broadcasting message: ", err)
				return err
			}
		}
	}
	return nil
}
