package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type AvatarsRepoI interface {
	InsertAvatar(avatar *models.AvatarDB) error
	UpdateAvatarByUserId(userId string, avatar *models.AvatarDB) error
	GetAvatarByUser(userId int64, avatarDB *models.AvatarDB) error
}
