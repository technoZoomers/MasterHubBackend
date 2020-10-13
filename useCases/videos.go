package useCases

import (
	"fmt"
	"github.com/google/logger"
	"github.com/h2non/filetype"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type VideosUC struct {
	VideosRepo  repository.VideosRepoI
	MastersRepo repository.MastersRepoI
}

func (videosUC *VideosUC) NewMasterVideo(videoData *models.VideoData, file multipart.File, masterId int64) (bool, error) {
	masterDB := models.MasterDB{
		UserId: masterId,
	}
	errType, err := videosUC.MastersRepo.GetMasterByUserId(&masterDB)
	if err != nil {
		if errType == utils.USER_ERROR {
			return true, err
		} else if errType == utils.SERVER_ERROR {
			return false, fmt.Errorf("database internal error")
		}
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fileError := fmt.Errorf("error reading file: %s", err.Error())
		logger.Errorf(fileError.Error())
		return false, fileError
	}

	fileExtension, err := filetype.Match(fileBytes)
	if err != nil {
		fileError := fmt.Errorf("error reading file extension: %s", err.Error())
		logger.Errorf(fileError.Error())
		return false, fileError
	}
	countVideo, err := videosUC.VideosRepo.CountVideos()
	if err != nil {
		return false, fmt.Errorf("database internal error")
	}
	newPath := filepath.Join(fmt.Sprintf("./master_videos/master_video_%d.%s", countVideo+1, fileExtension.Extension))
	newFile, err := os.Create(newPath)
	if err != nil {
		fileError := fmt.Errorf("error creating file: %s", err.Error())
		logger.Errorf(fileError.Error())
		return false, fileError
	}
	defer newFile.Close()

	_, err = newFile.Write(fileBytes)
	if err != nil {
		os.Remove(newPath)
		fileError := fmt.Errorf("error creating file: %s", err.Error())
		logger.Errorf(fileError.Error())
		return false, fileError
	}

	videoDB := models.VideoDB{
		Filename: newPath,
		MasterId: masterId,
		Name:     utils.DEFAULT_VIDEO_NAME,
		Intro:    false,
		Uploaded: time.Now(),
	}
	err = videosUC.VideosRepo.InsertVideoData(&videoDB)
	if err != nil {
		os.Remove(newPath)
		return false, fmt.Errorf("database internal error")
	}

	videoData.Name = videoDB.Name
	videoData.Uploaded = videoDB.Uploaded
	videoData.Id = videoDB.Id

	return false, nil
}

func (videosUC *VideosUC) GetVideosByMasterId(masterId int64) ([]models.VideoData, bool, error) {
	var videos []models.VideoData
	var videosDB []models.VideoDB
	masterDB := models.MasterDB{
		UserId: masterId,
	}
	errType, err := videosUC.MastersRepo.GetMasterByUserId(&masterDB)
	if err != nil {
		if errType == utils.USER_ERROR {
			return videos, true, err
		} else if errType == utils.SERVER_ERROR {
			return videos, false, fmt.Errorf("database internal error")
		}
	}
	videosDB, err = videosUC.VideosRepo.GetVideosByMasterId(masterId)
	for _, videoDB := range videosDB {
		video := models.VideoData{
			Id:          videoDB.Id,
			Name:        videoDB.Name,
			Description: videoDB.Description,
			Uploaded:    videoDB.Uploaded,
		}
		videos = append(videos, video)
	}
	return videos, false, nil
}