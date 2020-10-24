package useCases

import (
	"github.com/technoZoomers/MasterHubBackend/models"
)

type ChatsUCInterface interface {
	GetUserChatsById(userId int64, query models.ChatsQueryValues) (models.Chats, error)
	CreateChatRequest(chatRequest *models.Chat, studentId int64) error
}