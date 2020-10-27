package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type ChatsRepoI interface {
	GetChatsByUserId(query models.ChatsQueryValuesDB) ([]models.ChatDB, error)
	InsertChatRequest(chat *models.ChatDB) error
	GetChatByStudentIdAndMasterId(chat *models.ChatDB) error
	GetChatById(chat *models.ChatDB, chatId int64) error
	GetChatByIdAndMasterOrStudentId(chat *models.ChatDB, chatId int64, userId int64) error
	ChangeChatType(chat *models.ChatDB) error
	GetMessagesByChatId(chatId int64) ([]models.MessageDB, error)
	InsertMessage(message *models.MessageDB) error
}
