package useCases

import "github.com/technoZoomers/MasterHubBackend/models"

type UsersUCInterface interface {
	GetUserById(user *models.User) error
	Login(user *models.User) error
	GetUserByCookie(cookieValue string, user *models.User) error
	InsertCookie(userId int64, cookieValue string) error
	DeleteCookie(cookieValue string) error
}
