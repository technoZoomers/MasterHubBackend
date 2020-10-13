package useCases

import (
	"github.com/technoZoomers/MasterHubBackend/models"
	"mime/multipart"
	"os"
)

type VideosUCInterface interface {
	NewMasterVideo(videoData *models.VideoData, file multipart.File, id int64) (bool, error)
	GetVideosByMasterId(masterId int64) ([]models.VideoData, bool, error)
	GetMasterVideo(masterId int64, videoId int64) (*os.File, bool, error)
}
