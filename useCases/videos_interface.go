package useCases

import (
	"github.com/technoZoomers/MasterHubBackend/models"
	"mime/multipart"
)

type VideosUCInterface interface {
	NewMasterVideo(videoData *models.VideoData, file multipart.File, id int64) (bool, error)
	GetVideosByMasterId(masterId int64) ([]models.VideoData, bool, error)
	GetMasterVideo(masterId int64, videoId int64) ([]byte, bool, error)
	GetVideoDataById(videoData *models.VideoData, masterId int64) (bool, error)
	ChangeVideoData(videoData *models.VideoData, masterId int64) (bool, error)
}
