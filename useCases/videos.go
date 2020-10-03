package useCases

import (
	"fmt"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type VideosUC struct {
	VideosRepo  repository.VideosRepoI
}

func (videosUC *VideosUC) NewMasterVideo(videoData *models.VideoData, file multipart.File, masterId int64) (bool, error) {
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return false, err
	}
	fileType := http.DetectContentType(fileBytes)
	if fileType != "video/avi" && fileType != "video/webm" {
		fileTypeError := fmt.Errorf("wrong file format")
		return false, fileTypeError
	}
	fileEndings, err := mime.ExtensionsByType(fileType)
	if err != nil {
		return false, err
	}

	countVideo, err := videosUC.VideosRepo.CountVideos()
	countVideoString := strconv.FormatInt(countVideo+1, 10)

	newPath := filepath.Join("./master_videos/", "master_video_" + countVideoString +fileEndings[0])
	newFile, err := os.Create(newPath)
	if err != nil {
		return false, err
	}
	var videoDB models.VideoDB
	videoDB.Filename = newPath
	videoDB.MasterId = masterId
	videoDB.Intro = false
	videoDB.Uploaded = time.Now()
	err = videosUC.VideosRepo.InsertVideoData(&videoDB)
	if err != nil {
		os.Remove(newPath)
		return false, err
	}

	defer newFile.Close()
	_, err = newFile.Write(fileBytes)
	if  err != nil {
		return false, err
	}
	videoData.Name = utils.DEFAULT_VIDEO_NAME
	videoData.Uploaded = videoDB.Uploaded
	videoData.Id = videoDB.Id
	return false, nil

}