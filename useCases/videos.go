package useCases

import (
	"bufio"
	"fmt"
	"github.com/google/logger"
	"github.com/h2non/filetype"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/repository"
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
	videosDir         string
	videosDefaultName string
	videoPrefixMaster string
	videoPrefixVideo  string
	videoPrefixIntro  string
}

func (videosUC *VideosUC) createFilenameIntro(masterId int64) (string, error) {
	introExists := models.VideoDB{
		MasterId: masterId,
	}
	var filename string
	err := videosUC.VideosRepo.GetIntroByMasterId(&introExists)
	if err != nil {
		return filename, fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}
	if introExists.Id != videosUC.useCases.errorId {
		return filename, &models.ConflictError{Message: "intro already exists"}
	}
	filename = fmt.Sprintf("%s%d%s", videosUC.videosConfig.videoPrefixMaster, masterId, videosUC.videosConfig.videoPrefixIntro)
	return filename, nil
}

func (videosUC *VideosUC) createFilenameVideo(masterId int64) (string, error) {
	var filename string
	countVideo := videosUC.VideosRepo.GetLastVideoId()
	filename = fmt.Sprintf("%s%d%s%d", videosUC.videosConfig.videoPrefixMaster, masterId, videosUC.videosConfig.videoPrefixVideo, countVideo+1)
	return filename, nil
}

func (videosUC *VideosUC) validateMaster(masterId int64) (int64, error) {
	if masterId == videosUC.useCases.errorId {
		return videosUC.useCases.errorId, &models.BadRequestError{Message: "incorrect master id", RequestId: masterId}
	}
	masterDB := models.MasterDB{
		UserId: masterId,
	}
	err := videosUC.MastersRepo.GetMasterByUserId(&masterDB)
	if err != nil {
		return videosUC.useCases.errorId, fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}
	if masterDB.Id == videosUC.useCases.errorId {
		absenceError := &models.BadRequestError{Message: "master doesn't exist", RequestId: masterId}
		logger.Errorf(absenceError.Error())
		return videosUC.useCases.errorId, absenceError
	}
	return masterDB.Id, nil
}

func (videosUC *VideosUC) validateVideo(videoDB *models.VideoDB) error {
	if videoDB.Id == videosUC.useCases.errorId {
		return &models.BadRequestError{Message: "incorrect video id", RequestId: videoDB.Id}
	}

	err := videosUC.VideosRepo.GetVideoDataByIdAndMasterId(videoDB)
	if err != nil {
		return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}
	if videoDB.Name == "" {
		absenceError := &models.BadRequestError{Message: "video doesn't exist or doesn't belong to this master", RequestId: videoDB.Id}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	return nil
}
func (videosUC *VideosUC) validateIntro(intro *models.VideoDB) error {
	err := videosUC.VideosRepo.GetIntroByMasterId(intro)
	if err != nil {
		return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}
	if intro.Name == "" {
		absenceError := &models.NoContentError{Message: "no intro", RequestId: intro.MasterId}
		logger.Errorf(absenceError.Error())
		return absenceError
	}
	return nil
}

func (videosUC *VideosUC) newVideo(videoData *models.VideoData, file multipart.File, masterId int64, intro bool) error {
	masterDBId, err := videosUC.validateMaster(masterId)
	if err != nil {
		return err
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fileError := fmt.Errorf("%s: %s", videosUC.useCases.errorMessages.FileErrors.FileReadError, err.Error())
		logger.Errorf(fileError.Error())
		return fileError
	}
	defer file.Close()
	fileExtension, err := filetype.Match(fileBytes)
	if err != nil {
		fileError := fmt.Errorf("%s: %s", videosUC.useCases.errorMessages.FileErrors.FileReadExtensionError, err.Error())
		logger.Errorf(fileError.Error())
		return fileError
	}

	var filename string
	if intro {
		filename, err = videosUC.createFilenameIntro(masterId)
	} else {
		filename, err = videosUC.createFilenameVideo(masterId)
	}
	if err != nil {
		return err
	}
	newPath := fmt.Sprintf("%s%s%s.%s", videosUC.useCases.filesDir, videosUC.videosConfig.videosDir, filename, fileExtension.Extension)
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
		Filename:  filename,
		Extension: fileExtension.Extension,
		MasterId:  masterDBId,
		Name:      videosUC.videosConfig.videosDefaultName,
		Intro:     intro,
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
	videoData.Intro = intro
	videoData.FileExt = videoDB.Extension
	videoData.Rating = 0

	return nil
}

func (videosUC *VideosUC) NewMasterVideo(videoData *models.VideoData, file multipart.File, masterId int64) error {
	return videosUC.newVideo(videoData, file, masterId, false)
}
func (videosUC *VideosUC) NewMasterIntro(videoData *models.VideoData, file multipart.File, masterId int64) error {
	return videosUC.newVideo(videoData, file, masterId, true)
}

func (videosUC *VideosUC) ChangeMasterIntro(videoData *models.VideoData, file multipart.File, masterId int64) error {
	videoDB := models.VideoDB{
		Intro: true,
	}
	err := videosUC.deleteVideo(&videoDB, masterId)
	if err != nil {
		return err
	}
	err = videosUC.newVideo(videoData, file, masterId, true)
	if err != nil {
		return err
	}
	return nil
}

func (videosUC *VideosUC) matchVideo(videoDB *models.VideoDB, video *models.VideoData) error {
	video.Id = videoDB.Id
	video.MasterId = videoDB.MasterId
	video.Name = videoDB.Name
	video.FileExt = videoDB.Extension
	video.Description = videoDB.Description
	video.Uploaded = videoDB.Uploaded
	video.Intro = videoDB.Intro
	video.Rating = videoDB.Rating
	if videoDB.Theme != videosUC.useCases.errorId {
		err := videosUC.setTheme(video, videoDB.Theme)
		if err != nil {
			return err
		}
		err = videosUC.setSubThemes(video, videoDB)
		if err != nil {
			return err
		}
	}
	return nil
}

func (videosUC *VideosUC) GetVideosByMasterId(masterId int64) ([]models.VideoData, error) {
	var videos []models.VideoData
	var videosDB []models.VideoDB
	masterDBId, err := videosUC.validateMaster(masterId)
	if err != nil {
		return videos, err
	}
	videosDB, err = videosUC.VideosRepo.GetVideosByMasterId(masterDBId)
	if err != nil {
		return videos, fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}

	for _, videoDB := range videosDB {
		var video models.VideoData
		err = videosUC.matchVideo(&videoDB, &video)
		if err != nil {
			return videos, err
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
	//video.Theme.Id = theme
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

func (videosUC *VideosUC) deleteVideo(videoDB *models.VideoDB, masterId int64) error {
	masterDBId, err := videosUC.validateMaster(masterId)
	if err != nil {
		return err
	}
	videoDB.MasterId = masterDBId
	if videoDB.Intro {
		err = videosUC.validateIntro(videoDB)
	} else {
		err = videosUC.validateVideo(videoDB)
	}
	if err != nil {
		return err
	}
	err = videosUC.VideosRepo.DeleteVideo(videoDB)
	if err != nil {
		return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}
	filename := fmt.Sprintf("%s%s%s.%s", videosUC.useCases.filesDir, videosUC.videosConfig.videosDir, videoDB.Filename, videoDB.Extension)
	err = os.Remove(filename)
	if err != nil {
		fileError := fmt.Errorf("%s: %s", videosUC.useCases.errorMessages.FileErrors.FileRemoveError, err.Error())
		logger.Errorf(fileError.Error())
		return fileError
	}
	return nil
}

func (videosUC *VideosUC) DeleteMasterVideo(masterId int64, videoId int64) error {
	videoDB := models.VideoDB{
		Id:    videoId,
		Intro: false,
	}
	return videosUC.deleteVideo(&videoDB, masterId)
}

func (videosUC *VideosUC) DeleteMasterIntro(masterId int64) error {
	videoDB := models.VideoDB{
		Intro: true,
	}
	return videosUC.deleteVideo(&videoDB, masterId)
}

func (videosUC *VideosUC) getVideo(videoDB *models.VideoDB, masterId int64) ([]byte, error) {
	var videoBytes []byte
	masterDBId, err := videosUC.validateMaster(masterId)
	if err != nil {
		return videoBytes, err
	}
	videoDB.MasterId = masterDBId
	if videoDB.Intro {
		err = videosUC.validateIntro(videoDB)
	} else {
		err = videosUC.validateVideo(videoDB)
	}
	if err != nil {
		return videoBytes, err
	}
	filename := fmt.Sprintf("%s%s%s.%s", videosUC.useCases.filesDir, videosUC.videosConfig.videosDir, videoDB.Filename, videoDB.Extension)
	videoFile, err := os.Open(filename)
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

func (videosUC *VideosUC) GetMasterVideo(masterId int64, videoId int64) ([]byte, error) {
	videoDB := models.VideoDB{
		Id:    videoId,
		Intro: false,
	}
	return videosUC.getVideo(&videoDB, masterId)

}
func (videosUC *VideosUC) GetMasterIntro(masterId int64) ([]byte, error) {
	videoDB := models.VideoDB{
		Intro: true,
	}
	return videosUC.getVideo(&videoDB, masterId)
}

func (videosUC *VideosUC) getVideoData(videoDB *models.VideoDB, masterId int64, videoData *models.VideoData) error {
	masterDBId, err := videosUC.validateMaster(masterId)
	if err != nil {
		return err
	}
	videoDB.MasterId = masterDBId
	if videoDB.Intro {
		err = videosUC.validateIntro(videoDB)
	} else {
		err = videosUC.validateVideo(videoDB)
	}
	if err != nil {
		return err
	}
	err = videosUC.matchVideo(videoDB, videoData)
	if err != nil {
		return err
	}
	return nil
}

func (videosUC *VideosUC) GetVideoDataById(videoData *models.VideoData, masterId int64) error {
	videoDB := models.VideoDB{
		Id:    videoData.Id,
		Intro: false,
	}
	return videosUC.getVideoData(&videoDB, masterId, videoData)
}

func (videosUC *VideosUC) GetIntroData(videoData *models.VideoData, masterId int64) error {
	videoDB := models.VideoDB{
		Intro: true,
	}
	return videosUC.getVideoData(&videoDB, masterId, videoData)
}

func (videosUC *VideosUC) changeVideoData(videoDB *models.VideoDB, masterId int64, videoData *models.VideoData) error {
	masterDBId, err := videosUC.validateMaster(masterId)
	if err != nil {
		return err
	}
	videoDB.MasterId = masterDBId
	if videoDB.Intro {
		err = videosUC.validateIntro(videoDB)
	} else {
		err = videosUC.validateVideo(videoDB)
	}
	if err != nil {
		return err
	}
	if videoData.Rating != videosUC.useCases.errorId {
		if videoData.Rating != videoDB.Rating {
			requestError := fmt.Errorf("video rating can't be changed") //TODO: refactor error type
			logger.Errorf(requestError.Error())
			return requestError
		}
	} else {
		videoData.Rating = videoDB.Rating
	}
	if videoData.FileExt != "" {
		if videoData.FileExt != videoDB.Extension {
			requestError := fmt.Errorf("video extension can't be changed")
			logger.Errorf(requestError.Error())
			return requestError
		}
	} else {
		videoData.FileExt = videoDB.Extension
	}
	if !videoData.Uploaded.IsZero() {
		if !videoDB.Uploaded.Equal(videoData.Uploaded) {
			requestError := fmt.Errorf("video upload time can't be changed")
			logger.Errorf(requestError.Error())
			return requestError
		}
	} else {
		videoData.Uploaded = videoDB.Uploaded
	}

	if videoData.Name != "" {
		videoDB.Name = videoData.Name
	}
	videoDB.Description = videoData.Description

	err = videosUC.changeVideoTheme(videoData, videoDB)
	if err != nil {
		return err
	}

	err = videosUC.VideosRepo.UpdateVideo(videoDB)
	if err != nil {
		return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}

	err = videosUC.changeVideoSubthemes(videoData, videoDB)
	if err != nil {
		return err
	}
	return nil
}

func (videosUC *VideosUC) ChangeIntroData(videoData *models.VideoData, masterId int64) error {
	videoDB := models.VideoDB{
		Intro: true,
	}
	return videosUC.changeVideoData(&videoDB, masterId, videoData)
}

func (videosUC *VideosUC) ChangeVideoData(videoData *models.VideoData, masterId int64, videoId int64) error {
	if videoData.Id == videosUC.useCases.errorId {
		videoData.Id = videoId
	} else {
		if videoData.Id != videoId {
			matchError := &models.BadRequestError{Message: "video id doesn't match", RequestId: videoData.Id}
			logger.Errorf(matchError.Error())
			return matchError
		}
	}
	videoDB := models.VideoDB{
		Id:    videoData.Id,
		Intro: false,
	}
	return videosUC.changeVideoData(&videoDB, masterId, videoData)
}

func (videosUC *VideosUC) changeVideoSubthemes(videoData *models.VideoData, videoDB *models.VideoDB) error {
	var err error
	if videoData.Theme.Theme == "" {
		err = videosUC.MastersRepo.DeleteMasterSubthemesById(videoDB.Id)
		if err != nil {
			return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
		}
		return nil
	}

	var newSubthemesIds []int64
	for _, subtheme := range videoData.Theme.Subthemes {
		subthemeDB := models.SubthemeDB{Name: subtheme}
		err := videosUC.ThemesRepo.GetSubthemeByName(&subthemeDB)
		if err != nil {
			return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
		}
		if subthemeDB.Id == videosUC.useCases.errorId {
			fileError := &models.BadRequestError{Message: "cant't update video, subtheme doesn't exist", RequestId: videoData.Id}
			logger.Errorf(fileError.Error())
			return fileError
		}
		newSubthemesIds = append(newSubthemesIds, subthemeDB.Id)
	}

	oldSubthemesIds, err := videosUC.VideosRepo.GetVideoSubthemesById(videoDB.Id)
	if err != nil {
		return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}

	err = videosUC.VideosRepo.DeleteVideoSubthemesById(videoData.Id)
	if err != nil {
		return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}

	err = videosUC.VideosRepo.SetVideoSubthemesById(videoData.Id, newSubthemesIds)
	if err != nil {
		_ = videosUC.VideosRepo.SetVideoSubthemesById(videoData.Id, oldSubthemesIds)
		return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}
	return nil
}

func (videosUC *VideosUC) changeVideoTheme(videoData *models.VideoData, videoDB *models.VideoDB) error {

	if videoData.Theme.Theme == "" {
		videoDB.Theme = videosUC.useCases.errorId
		return nil
	}

	var oldTheme models.ThemeDB
	oldTheme.Id = videoDB.Theme
	err := videosUC.getTheme(&oldTheme)
	if err != nil {
		return err
	}

	if videoData.Theme.Theme != oldTheme.Name {
		newThemeDB := models.ThemeDB{
			Name: videoData.Theme.Theme,
		}
		err := videosUC.ThemesRepo.GetThemeByName(&newThemeDB)
		if err != nil {
			return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
		}
		if newThemeDB.Id == videosUC.useCases.errorId {
			fileError := &models.BadRequestError{Message: "cant't update video, theme doesn't exist", RequestId: videoData.Id}
			logger.Errorf(fileError.Error())
			return fileError
		}
		videoDB.Theme = newThemeDB.Id
	}
	return nil
}

func (videosUC *VideosUC) matchTheme(theme string, queryDB *models.VideosQueryValuesDB) error {
	if theme != "" {
		themeDB := models.ThemeDB{
			Name: theme,
		}
		err := videosUC.ThemesRepo.GetThemeByName(&themeDB)
		if err != nil {
			return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
		}
		if themeDB.Id == videosUC.useCases.errorId {
			badParamError := &models.BadQueryParameterError{Parameter: "theme"}
			logger.Errorf(badParamError.Error())
			return badParamError
		}
		queryDB.Theme = append(queryDB.Theme, themeDB.Id)
	}
	return nil
}

func (videosUC *VideosUC) matchSubthemes(subthemes []string, queryDB *models.VideosQueryValuesDB) error {
	for _, subtheme := range subthemes {
		subthemeDB := models.SubthemeDB{Name: subtheme}
		err := videosUC.ThemesRepo.GetSubthemeByName(&subthemeDB)
		if err != nil {
			return fmt.Errorf(videosUC.useCases.errorMessages.DbError)
		}
		if subthemeDB.Id == videosUC.useCases.errorId {
			badParamError := &models.BadQueryParameterError{Parameter: "subtheme"}
			logger.Errorf(badParamError.Error())
			return badParamError
		}
		queryDB.Subtheme = append(queryDB.Subtheme, subthemeDB.Id)
	}
	return nil
}

func (videosUC *VideosUC) matchVideosQuery(query *models.VideosQueryValues, queryDB *models.VideosQueryValuesDB) error {
	queryDB.Offset = query.Offset
	queryDB.Limit = query.Limit
	queryDB.Old = query.Old
	queryDB.Popular = query.Popular
	err := videosUC.matchTheme(query.Theme, queryDB)
	if err != nil {
		return err
	}
	queryDB.Subtheme = make([]int64, 0)
	err = videosUC.matchSubthemes(query.Subtheme, queryDB)
	if err != nil {
		return err
	}
	return nil
}

func (videosUC *VideosUC) Get(query models.VideosQueryValues) (models.VideosData, error) {
	var queryDB models.VideosQueryValuesDB
	videos := make([]models.VideoData, 0)
	err := videosUC.matchVideosQuery(&query, &queryDB)
	if err != nil {
		return videos, err
	}
	videosDB, err := videosUC.VideosRepo.GetVideos(queryDB)
	if err != nil {
		return videos, fmt.Errorf(videosUC.useCases.errorMessages.DbError)
	}
	for _, videoDB := range videosDB {
		var videoData models.VideoData
		err = videosUC.matchVideo(&videoDB, &videoData)
		if err != nil {
			return videos, err
		}
		videos = append(videos, videoData)
	}
	return videos, nil
}
