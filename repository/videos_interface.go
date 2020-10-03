package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type VideosRepoI interface {
	InsertVideoData(video *models.VideoDB) error
	CountVideos() (int64, error)
}
