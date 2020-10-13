package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type VideosRepoI interface {
	InsertVideoData(video *models.VideoDB) error
	CountVideos() (int64, error)
	GetVideosByMasterId(masterId int64) ([]models.VideoDB, error)
}
