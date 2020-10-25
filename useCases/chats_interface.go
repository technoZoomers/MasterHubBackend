package useCases

import (
	"github.com/technoZoomers/MasterHubBackend/models"
)

type ChatsUCInterface interface {
	GetUserChatsById(userId int64, query models.ChatsQueryValues) (models.Chats, error)
	CreateChatRequest(chatRequest *models.Chat, studentId int64) error
	ChangeChatStatus(chat *models.Chat, masterId int64, chatId int64) error
}