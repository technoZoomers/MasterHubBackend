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
	VideosRepo  repository.VideosRepoI
	MastersRepo repository.MastersRepoI
	ThemesRepo repository.ThemesRepoI
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
	fileName := fmt.Sprintf("master_video_%d", countVideo+1)
	newPath := fmt.Sprintf("./master_videos/%s.%s", fileName, fileExtension.Extension)
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
		Filename: fileName,
		Extension: fileExtension.Extension,
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
	videoData.FileExt = videoDB.Extension

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
	if err != nil {
		return  videos, false, fmt.Errorf("database internal error")
	}

	for _, videoDB := range videosDB {
		video := models.VideoData{
			Id:          videoDB.Id,
			Name:        videoDB.Name,
			FileExt: videoDB.Extension,
			Description: videoDB.Description,
			Uploaded:    videoDB.Uploaded,
		}
		if videoDB.Theme != 0 {
			err = videosUC.setTheme(&video, videoDB.Theme)
			if err != nil {
				return videos, false, err
			}
			err = videosUC.setSubThemes(&video, &videoDB)
			if err != nil {
				return videos, false, err
			}
		}

		videos = append(videos, video)
	}
	return videos, false, nil
}


func (videosUC *VideosUC) getTheme(themeDB *models.ThemeDB) error {
	_, err := videosUC.ThemesRepo.GetThemeById(themeDB)
	if err != nil {
		return err
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
		return err
	}
	for _, subthemeId := range subthemesIds {
		var subtheme models.SubthemeDB
		subtheme.Id = subthemeId
		_, err = videosUC.ThemesRepo.GetSubthemeById(&subtheme)
		if err != nil {
			video.Theme.Subthemes = subthemes
			return err
		}
		subthemes = append(subthemes, subtheme.Name)
	}
	video.Theme.Subthemes = subthemes
	return nil
}

func (videosUC *VideosUC) GetMasterVideo(masterId int64, videoId int64) ([]byte, bool, error) {
	var videoBytes []byte
	videoDB := models.VideoDB{
		Id:          videoId,
		MasterId: masterId,
	}
	masterDB := models.MasterDB{
		UserId: masterId,
	}
	errType, err := videosUC.MastersRepo.GetMasterByUserId(&masterDB)
	if err != nil {
		if errType == utils.USER_ERROR {
			return videoBytes, true, err
		} else if errType == utils.SERVER_ERROR {
			return videoBytes, false, fmt.Errorf("database internal error")
		}
	}
	errType, err = videosUC.VideosRepo.GetVideoDataById(&videoDB)
	if err != nil {
		if errType == utils.USER_ERROR {
			return videoBytes, true, err
		} else if errType == utils.SERVER_ERROR {
			return videoBytes, false, fmt.Errorf("database internal error")
		}
	}
	videoFile, err := os.Open(fmt.Sprintf("./master_videos/%s.%s", videoDB.Filename, videoDB.Extension))
	if err != nil {
		fileError := fmt.Errorf("error opening file: %s", err.Error())
		logger.Errorf(fileError.Error())
		return videoBytes, false, fileError
	}
	defer videoFile.Close()

	reader := bufio.NewReader(videoFile)
	videoFileInfo, err := videoFile.Stat()
	if err != nil {
		fileError := fmt.Errorf("error opening file: %s", err.Error())
		logger.Errorf(fileError.Error())
		return videoBytes, false, fileError
	}
	videoFileSize := videoFileInfo.Size()

	videoBytes = make([]byte, videoFileSize)
	_, err = reader.Read(videoBytes)
	if err != nil {
		fileError := fmt.Errorf("error reading file: %s", err.Error())
		logger.Errorf(fileError.Error())
		return videoBytes, false, fileError
	}
	return videoBytes, false, nil
}

func (videosUC *VideosUC) GetVideoDataById(videoData *models.VideoData, masterId int64) (bool, error) {
	videoDB := models.VideoDB{
		Id: videoData.Id,
		MasterId: masterId,
	}
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
	errType, err = videosUC.VideosRepo.GetVideoDataById(&videoDB)
	if err != nil {
		if errType == utils.USER_ERROR {
			return true, err
		} else if errType == utils.SERVER_ERROR {
			return false, fmt.Errorf("database internal error")
		}
	}

	videoData.Name = videoDB.Name
	videoData.FileExt = videoDB.Extension
	videoData.Description = videoDB.Description
	videoData.Uploaded = videoDB.Uploaded

	if videoDB.Theme != 0 {
		err = videosUC.setTheme(videoData, videoDB.Theme)
		if err != nil {
			return false, err
		}
		err = videosUC.setSubThemes(videoData, &videoDB)
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

func (videosUC *VideosUC) ChangeVideoData(videoData *models.VideoData, masterId int64) (bool, error) {
	videoDB := models.VideoDB{
		Id: videoData.Id,
		MasterId: masterId,
	}
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
	errType, err = videosUC.VideosRepo.GetVideoDataById(&videoDB)
	if err != nil {
		if errType == utils.USER_ERROR {
			return true, err
		} else if errType == utils.SERVER_ERROR {
			return false, fmt.Errorf("database internal error")
		}
	}

	if videoData.FileExt != videoDB.Extension {
		fileError := fmt.Errorf("video extension can't be changed")
		logger.Errorf(fileError.Error())
		return false, fileError
	}
	fmt.Println(videoDB.Uploaded, videoData.Uploaded)
	fmt.Println(videoDB.Uploaded.Equal(videoData.Uploaded))
	if videoData.Uploaded != videoDB.Uploaded {
		fileError := fmt.Errorf("video upload time can't be changed")
		logger.Errorf(fileError.Error())
		return false, fileError
	}
	var themeDB models.ThemeDB
	themeDB.Id = videoDB.Theme
	_ = videosUC.getTheme(&themeDB)

	absent, err := videosUC.changeVideoTheme(videoData, &themeDB, &videoDB)
	if err != nil {
		return absent, err
	}
	videoDB.Description = videoData.Description
	videoDB.Name = videoData.Name
	err = videosUC.VideosRepo.UpdateVideo(&videoDB)
	if err != nil {
		return false, fmt.Errorf("database internal error")
	}
	return false, nil
}

func (videosUC *VideosUC) changeVideoTheme(videoData *models.VideoData, oldTheme *models.ThemeDB, videoDB *models.VideoDB) (bool, error) {
	if videoData.Theme.Theme != oldTheme.Name {
		newThemeDB := models.ThemeDB{
			Name:videoData.Theme.Theme,
		}
		errType, err := videosUC.ThemesRepo.GetThemeByName(&newThemeDB)
		if err != nil {
			if errType == utils.USER_ERROR {
				fileError := fmt.Errorf("cant't update video, theme doesn't exist")
				logger.Errorf(fileError.Error())
				return true, fileError
			} else if errType == utils.SERVER_ERROR {
				return false, fmt.Errorf("database internal error")
			}
		}
		videoDB.Theme = newThemeDB.Id
	}
	var newSubthemesIds []int64
	for _, subtheme := range videoData.Theme.Subthemes {
		subthemeDB := models.SubthemeDB{Name:subtheme}
		errType, err := videosUC.ThemesRepo.GetSubthemeByName(&subthemeDB)
		if err != nil {
			if errType == utils.USER_ERROR {
				fileError := fmt.Errorf("cant't update video, subtheme doesn't exist")
				logger.Errorf(fileError.Error())
				return true, fileError
			} else if errType == utils.SERVER_ERROR {
				return false, fmt.Errorf("database internal error")
			}
		}
		newSubthemesIds = append(newSubthemesIds, subthemeDB.Id)
	}

	err := videosUC.VideosRepo.DeleteVideoSubthemesById(videoData.Id)
	if err != nil {
		return false, fmt.Errorf("database internal error")
	}

	err = videosUC.VideosRepo.SetVideoSubthemesById(videoData.Id, newSubthemesIds)
	if err != nil {
		videoData.Theme.Subthemes = []string{}
		return false, fmt.Errorf("database internal error")
	}
	return false, nil
}
