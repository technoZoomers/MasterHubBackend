package useCases

import "github.com/technoZoomers/MasterHubBackend/models"

type UsersUCInterface interface {
	GetUserById(user *models.User) error
	Login(user *models.User) error

}
