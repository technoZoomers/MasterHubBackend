package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type UsersRepoI interface {
	InsertUser(user *models.UserDB) error
	GetUserById(user *models.UserDB, userId int64) error
	GetUserByEmail(user *models.UserDB) error
	GetUserByEmailAndPassword(user *models.UserDB) error
	GetUserLanguagesById(userId int64) ([]int64, error)
	DeleteUserLanguagesById(userId int64) error
	SetUserLanguagesById(userId int64, languages []int64) error
	DeleteUserWithId(userId int64) error
}
