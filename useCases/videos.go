package useCases

import (
	"bufio"
	"fmt"
	"github.com/google/logger"
	"github.com/h2non/filetype"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"io/ioutil"
	"mime/multipart"
	"os"
	"time"
)

type VideosUC struct {
	useCases     *UseCases
	VideosRepo   repository.VideosRepoI
	MastersRepo  repository.MastersRepoI
	ThemesRepo   repository.ThemesRepoI
	videosConfig VideoConfig
}

type VideoConfig struct {
	videosDir           string
	videosDefaultName   string
	videoFilenamePrefix string
}

func (videosUC *VideosUC) NewMasterVideo(videoData *models.VideoData, file multipart.File, masterId int64) error {
	if masterId == utils.ERROR_ID {
		return &models.BadRequestError{Message: "incorrect master id", RequestId: masterId}
	}
	masterDB := models.MasterDB{
		UserId: masterId,
	}
	err := videosUC.MastersRepo.GetMasterByUserId(&masterDB)
	if err != nil {
		return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}
	if masterDB.Id == utils.ERROR_ID {
		absenceError := &models.BadRequestError{Message: "master doesn't exist", RequestId: masterId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fileError := fmt.Errorf("%s: %s", videosUC.useCases.errorMessages.FileErrors.FileReadError, err.Error())
		logger.Errorf(fileError.Error())
		return fileError
	}

	fileExtension, err := filetype.Match(fileBytes)
	if err != nil {
		fileError := fmt.Errorf("%s: %s", videosUC.useCases.errorMessages.FileErrors.FileReadExtensionError, err.Error())
		logger.Errorf(fileError.Error())
		return fileError
	}
	countVideo, err := videosUC.VideosRepo.CountVideos()
	if err != nil {
		return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}
	fileName := fmt.Sprintf("%s%d", videosUC.videosConfig.videoFilenamePrefix, countVideo+1)
	newPath := fmt.Sprintf("%s%s.%s", videosUC.videosConfig.videosDir, fileName, fileExtension.Extension)
	newFile, err := os.Create(newPath)
	if err != nil {
		fileError := fmt.Errorf("%s: %s", videosUC.useCases.errorMessages.FileErrors.FileCreateError, err.Error())
		logger.Errorf(fileError.Error())
		return fileError
	}
	defer newFile.Close()

	_, err = newFile.Write(fileBytes)
	if err != nil {
		os.Remove(newPath)
		fileError := fmt.Errorf("%s: %s", videosUC.useCases.errorMessages.FileErrors.FileCreateError, err.Error())
		logger.Errorf(fileError.Error())
		return fileError
	}

	videoDB := models.VideoDB{
		Filename:  fileName,
		Extension: fileExtension.Extension,
		MasterId:  masterId,
		Name:      videosUC.videosConfig.videosDefaultName,
		Intro:     false,
		Uploaded:  time.Now(),
	}
	err = videosUC.VideosRepo.InsertVideoData(&videoDB)
	if err != nil {
		os.Remove(newPath)
		return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}

	videoData.Name = videoDB.Name
	videoData.Uploaded = videoDB.Uploaded
	videoData.Id = videoDB.Id
	videoData.FileExt = videoDB.Extension

	return nil
}

func (videosUC *VideosUC) GetVideosByMasterId(masterId int64) ([]models.VideoData, error) {
	var videos []models.VideoData
	var videosDB []models.VideoDB
	if masterId == utils.ERROR_ID {
		return videos, &models.BadRequestError{Message: "incorrect master id", RequestId: masterId}
	}
	masterDB := models.MasterDB{
		UserId: masterId,
	}
	err := videosUC.MastersRepo.GetMasterByUserId(&masterDB)
	if err != nil {
		return videos, fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}
	if masterDB.Id == utils.ERROR_ID {
		absenceError := &models.BadRequestError{Message: "master doesn't exist", RequestId: masterId}
		logger.Errorf(absenceError.Error())
		return videos, absenceError
	}
	videosDB, err = videosUC.VideosRepo.GetVideosByMasterId(masterId)
	if err != nil {
		return videos, fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}

	for _, videoDB := range videosDB {
		video := models.VideoData{
			Id:          videoDB.Id,
			Name:        videoDB.Name,
			FileExt:     videoDB.Extension,
			Description: videoDB.Description,
			Uploaded:    videoDB.Uploaded,
		}
		if videoDB.Theme != utils.ERROR_ID {
			err = videosUC.setTheme(&video, videoDB.Theme)
			if err != nil {
				return videos, err
			}
			err = videosUC.setSubThemes(&video, &videoDB)
			if err != nil {
				return videos, err
			}
		}

		videos = append(videos, video)
	}
	return videos, nil
}

func (videosUC *VideosUC) getTheme(themeDB *models.ThemeDB) error {
	err := videosUC.ThemesRepo.GetThemeById(themeDB)
	if err != nil {
		return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}
	return nil
}

func (videosUC *VideosUC) setTheme(video *models.VideoData, theme int64) error {
	var themeDB models.ThemeDB
	themeDB.Id = theme
	err := videosUC.getTheme(&themeDB)
	if err != nil {
		return err
	}
	video.Theme.Id = theme
	video.Theme.Theme = themeDB.Name
	return nil
}

func (videosUC *VideosUC) setSubThemes(video *models.VideoData, videoDB *models.VideoDB) error {
	var subthemes []string
	subthemesIds, err := videosUC.VideosRepo.GetVideoSubthemesById(videoDB.Id)
	if err != nil {
		video.Theme.Subthemes = subthemes
		return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}
	for _, subthemeId := range subthemesIds {
		var subtheme models.SubthemeDB
		subtheme.Id = subthemeId
		err = videosUC.ThemesRepo.GetSubthemeById(&subtheme)
		if err != nil {
			video.Theme.Subthemes = subthemes
			return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
		}
		if subtheme.Name == "" {
			absenceError := fmt.Errorf("subtheme doesn't exist")
			logger.Errorf(absenceError.Error())
			video.Theme.Subthemes = subthemes
			return absenceError
		}
		subthemes = append(subthemes, subtheme.Name)
	}
	video.Theme.Subthemes = subthemes
	return nil
}

func (videosUC *VideosUC) GetMasterVideo(masterId int64, videoId int64) ([]byte, error) {
	var videoBytes []byte
	if masterId == utils.ERROR_ID {
		return videoBytes, &models.BadRequestError{Message: "incorrect master id", RequestId: masterId}
	}
	if videoId == utils.ERROR_ID {
		return videoBytes, &models.BadRequestError{Message: "incorrect video id", RequestId: videoId}
	}
	videoDB := models.VideoDB{
		Id:       videoId,
		MasterId: masterId,
	}
	masterDB := models.MasterDB{
		UserId: masterId,
	}
	err := videosUC.MastersRepo.GetMasterByUserId(&masterDB)
	if err != nil {
		return videoBytes, err
	}
	if masterDB.Id == utils.ERROR_ID {
		absenceError := &models.BadRequestError{Message: "master doesn't exist", RequestId: masterId}
		logger.Errorf(absenceError.Error())
		return videoBytes, absenceError
	}
	err = videosUC.VideosRepo.GetVideoDataById(&videoDB)
	if err != nil {
		return videoBytes, fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}
	if videoDB.Name == "" {
		absenceError := &models.BadRequestError{Message: "video doesn't exist or doesn't belong to this master", RequestId: videoId}
		logger.Errorf(absenceError.Error())
		return videoBytes, absenceError
	}
	videoFile, err := os.Open(fmt.Sprintf("%s%s.%s", videosUC.videosConfig.videosDir, videoDB.Filename, videoDB.Extension))
	if err != nil {
		fileError := fmt.Errorf("%s: %s", videosUC.useCases.errorMessages.FileErrors.FileOpenError, err.Error())
		logger.Errorf(fileError.Error())
		return videoBytes, fileError
	}
	defer videoFile.Close()

	reader := bufio.NewReader(videoFile)
	videoFileInfo, err := videoFile.Stat()
	if err != nil {
		fileError := fmt.Errorf("%s: %s", videosUC.useCases.errorMessages.FileErrors.FileOpenError, err.Error())
		logger.Errorf(fileError.Error())
		return videoBytes, fileError
	}
	videoFileSize := videoFileInfo.Size()

	videoBytes = make([]byte, videoFileSize)
	_, err = reader.Read(videoBytes)
	if err != nil {
		fileError := fmt.Errorf("%s: %s", videosUC.useCases.errorMessages.FileErrors.FileReadError, err.Error())
		logger.Errorf(fileError.Error())
		return videoBytes, fileError
	}
	return videoBytes, nil
}

func (videosUC *VideosUC) GetVideoDataById(videoData *models.VideoData, masterId int64) error {
	if masterId == utils.ERROR_ID {
		return &models.BadRequestError{Message: "incorrect master id", RequestId: masterId}
	}
	if videoData.Id == utils.ERROR_ID {
		return &models.BadRequestError{Message: "incorrect video id", RequestId: videoData.Id}
	}
	videoDB := models.VideoDB{
		Id:       videoData.Id,
		MasterId: masterId,
	}
	masterDB := models.MasterDB{
		UserId: masterId,
	}
	err := videosUC.MastersRepo.GetMasterByUserId(&masterDB)
	if err != nil {
		return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}
	if masterDB.Id == utils.ERROR_ID {
		absenceError := &models.BadRequestError{Message: "master doesn't exist", RequestId: masterId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	err = videosUC.VideosRepo.GetVideoDataById(&videoDB)
	if err != nil {
		return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}
	if videoDB.Name == "" {
		absenceError := &models.BadRequestError{Message: "video doesn't exist or doesn't belong to this master", RequestId: videoData.Id}
		logger.Errorf(absenceError.Error())
		return absenceError
	}

	videoData.Name = videoDB.Name
	videoData.FileExt = videoDB.Extension
	videoData.Description = videoDB.Description
	videoData.Uploaded = videoDB.Uploaded

	if videoDB.Theme != 0 {
		err = videosUC.setTheme(videoData, videoDB.Theme)
		if err != nil {
			return err
		}
		err = videosUC.setSubThemes(videoData, &videoDB)
		if err != nil {
			return err
		}
	}
	return nil
}

func (videosUC *VideosUC) ChangeVideoData(videoData *models.VideoData, masterId int64) error {
	if masterId == utils.ERROR_ID {
		return &models.BadRequestError{Message: "incorrect master id", RequestId: masterId}
	}
	if videoData.Id == utils.ERROR_ID {
		return &models.BadRequestError{Message: "incorrect video id", RequestId: videoData.Id}
	}
	videoDB := models.VideoDB{
		Id:       videoData.Id,
		MasterId: masterId,
	}
	masterDB := models.MasterDB{
		UserId: masterId,
	}
	err := videosUC.MastersRepo.GetMasterByUserId(&masterDB)
	if err != nil {
		return err
	}
	if masterDB.Id == utils.ERROR_ID {
		absenceError := &models.BadRequestError{Message: "master doesn't exist", RequestId: masterId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	err = videosUC.VideosRepo.GetVideoDataById(&videoDB)
	if err != nil {
		return fmt.Errorf("database internal error")
	}
	if videoDB.Name == "" {
		absenceError := &models.BadRequestError{Message: "video doesn't exist or doesn't belong to master", RequestId: videoData.Id}
		logger.Errorf(absenceError.Error())
		return absenceError
	}

	if videoData.FileExt != "" && videoData.FileExt != videoDB.Extension {
		fileError := fmt.Errorf("video extension can't be changed")
		logger.Errorf(fileError.Error())
		return fileError
	}
	if !videoData.Uploaded.IsZero() && !videoDB.Uploaded.Equal(videoData.Uploaded) {
		fileError := fmt.Errorf("video upload time can't be changed")
		logger.Errorf(fileError.Error())
		return fileError
	}
	var themeDB models.ThemeDB
	themeDB.Id = videoDB.Theme
	_ = videosUC.getTheme(&themeDB)

	err = videosUC.changeVideoTheme(videoData, &themeDB, &videoDB)
	if err != nil {
		return err
	}
	videoDB.Description = videoData.Description
	videoDB.Name = videoData.Name
	err = videosUC.VideosRepo.UpdateVideo(&videoDB)
	if err != nil {
		return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}
	return nil
}

func (videosUC *VideosUC) changeVideoTheme(videoData *models.VideoData, oldTheme *models.ThemeDB, videoDB *models.VideoDB) error {
	if videoData.Theme.Theme != oldTheme.Name {
		newThemeDB := models.ThemeDB{
			Name: videoData.Theme.Theme,
		}
		err := videosUC.ThemesRepo.GetThemeByName(&newThemeDB)
		if err != nil {
			return fmt.Errorf("database internal error")
		}
		if newThemeDB.Id == utils.ERROR_ID {
			fileError := &models.BadRequestError{Message: "cant't update video, theme doesn't exist", RequestId: videoData.Id}
			logger.Errorf(fileError.Error())
			return fileError
		}
		videoDB.Theme = newThemeDB.Id
	}
	var newSubthemesIds []int64
	for _, subtheme := range videoData.Theme.Subthemes {
		subthemeDB := models.SubthemeDB{Name: subtheme}
		err := videosUC.ThemesRepo.GetSubthemeByName(&subthemeDB)
		if err != nil {
			return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
		}
		if subthemeDB.Id == utils.ERROR_ID {
			fileError := &models.BadRequestError{Message: "cant't update video, subtheme doesn't exist", RequestId: videoData.Id}
			logger.Errorf(fileError.Error())
			return fileError
		}
		newSubthemesIds = append(newSubthemesIds, subthemeDB.Id)
	}

	err := videosUC.VideosRepo.DeleteVideoSubthemesById(videoData.Id)
	if err != nil {
		return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}

	err = videosUC.VideosRepo.SetVideoSubthemesById(videoData.Id, newSubthemesIds)
	if err != nil {
		videoData.Theme.Subthemes = []string{}
		return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}
	return nil
}
