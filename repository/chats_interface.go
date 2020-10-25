package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type ChatsRepoI interface {
	GetChatsByUserId(query models.ChatsQueryValuesDB) ([]models.ChatDB, error)
	InsertChatRequest(chat *models.ChatDB) error
	GetChatByStudentIdAndMasterId(chat *models.ChatDB) error
	GetChatById(chat *models.ChatDB, chatId int64) error
	ChangeChatType(chat *models.ChatDB) error
}