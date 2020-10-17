package useCases

import (
	"github.com/technoZoomers/MasterHubBackend/models"
	"mime/multipart"
)

type VideosUCInterface interface {
	NewMasterVideo(videoData *models.VideoData, file multipart.File, masterId int64) error
	NewMasterIntro(videoData *models.VideoData, file multipart.File, masterId int64) error
	GetVideosByMasterId(masterId int64) ([]models.VideoData, error)
	GetMasterVideo(masterId int64, videoId int64) ([]byte, error)
	DeleteMasterVideo(masterId int64, videoId int64) error
	GetMasterIntro(masterId int64) ([]byte, error)
	GetVideoDataById(videoData *models.VideoData, masterId int64) error
	ChangeVideoData(videoData *models.VideoData, masterId int64) error
}
