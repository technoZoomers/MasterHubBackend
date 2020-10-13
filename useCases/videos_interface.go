package useCases

import (
	"github.com/technoZoomers/MasterHubBackend/models"
	"mime/multipart"
)

type VideosUCInterface interface {
	NewMasterVideo(videoData *models.VideoData, file multipart.File, id int64) (bool, error)
}
