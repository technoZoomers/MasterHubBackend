package useCases

import (
	"mime/multipart"
)

type AvatarsUCInterface interface {
	NewUserAvatar(file multipart.File, userId int64) error
	ChangeUserAvatar(file multipart.File, userId int64) error
	GetUserAvatar(userId int64) ([]byte, error)
}
