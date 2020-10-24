package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type ChatsRepoI interface {
	GetChatsByUserId(query models.ChatsQueryValuesDB) ([]models.ChatDB, error)
	InsertChatRequest(chat *models.ChatDB) error
}
