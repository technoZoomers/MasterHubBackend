package repository

import "github.com/technoZoomers/MasterHubBackend/models"

type VideosRepoI interface {
	InsertVideoData(video *models.VideoDB) error
	CountVideos() (int64, error)
	GetVideosByMasterId(masterId int64) ([]models.VideoDB, error)
	GetVideoDataByIdAndMasterId(video *models.VideoDB) error
	GetVideoSubthemesById(videoId int64) ([]int64, error)
	DeleteVideoSubthemesById(videoId int64) error
	SetVideoSubthemesById(videoId int64, subthemes []int64) error
	UpdateVideo(video *models.VideoDB) error
	DeleteVideo(video *models.VideoDB) error
	GetIntroByMasterId(video *models.VideoDB) error
}
