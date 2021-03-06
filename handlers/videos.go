package handlers

import (
	"fmt"
	"github.com/google/logger"
	json "github.com/mailru/easyjson"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"mime/multipart"
	"net/http"
)

type VideosHandlers struct {
	handlers     *Handlers
	VideosUC useCases.VideosUCInterface
	VideoParseConfig VideoParseConfig
}

type VideoParseConfig struct {
	FormDataKey string
	VideoFormats map[string]bool
}

func (vh *VideosHandlers) validateMasterId(writer http.ResponseWriter, req *http.Request) (bool, int64) {
	return vh.handlers.validateId(writer, req, "id", "master")
}

func (vh *VideosHandlers) validateVideoId(writer http.ResponseWriter, req *http.Request) (bool, int64) {
	return vh.handlers.validateId(writer, req, "videoId", "video")
}

func (vh *VideosHandlers) Upload(writer http.ResponseWriter, req *http.Request) {
	vh.uploadVideo(writer, req, false)
}

func (vh *VideosHandlers) UploadIntro(writer http.ResponseWriter, req *http.Request) {
	vh.uploadVideo(writer, req, true)
}

func (vh *VideosHandlers) getFileFromFormData (writer http.ResponseWriter, req *http.Request) (bool, multipart.File, error) {
	file, fileHeader, err := req.FormFile(vh.VideoParseConfig.FormDataKey)
	if err != nil {
		parseError := fmt.Errorf("error parsing video: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(parseError.Error()))
		return true, file, err
	}
	if !vh.VideoParseConfig.VideoFormats[fileHeader.Header.Get("Content-Type")] {
		parseError := fmt.Errorf("wrong mime type:%s, expected video", fileHeader.Header.Get("Content-Type"))
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(parseError.Error()))
		return true, file, err
	}
	return false, file, err
}

func (vh *VideosHandlers) uploadVideo (writer http.ResponseWriter, req *http.Request, intro bool) {
	var err error
	var videoData models.VideoData
	sent, masterId := vh.validateMasterId(writer, req)
	if sent {
		return
	}
	sent, file, err := vh.getFileFromFormData(writer, req)
	if sent {
		return
	}
	if intro {
		err = vh.VideosUC.NewMasterIntro(&videoData, file, masterId)
		vh.answerIntroPost(writer, videoData, http.StatusCreated, err)
	} else {
		err = vh.VideosUC.NewMasterVideo(&videoData, file, masterId)
		vh.answerVideo(writer, videoData, http.StatusCreated, err)
	}
}

func (vh *VideosHandlers) GetVideosByMasterId(writer http.ResponseWriter, req *http.Request) {
	var err error
	sent, masterId := vh.validateMasterId(writer, req)
	if sent {
		return
	}
	videos, err := vh.VideosUC.GetVideosByMasterId(masterId)
	vh.answerVideos(writer, videos, err)
}

func (vh *VideosHandlers) getVideo(writer http.ResponseWriter, req *http.Request, intro bool) {
	var err error
	sent, masterId := vh.validateMasterId(writer, req)
	if sent {
		return
	}
	var videoBytes []byte
	if intro {
		videoBytes, err = vh.VideosUC.GetMasterIntro(masterId)
		vh.answerMultipartIntro(writer, videoBytes, http.StatusOK, err)
	} else {
		sent, videoId := vh.validateVideoId(writer, req)
		if sent {
			return
		}
		videoBytes, err = vh.VideosUC.GetMasterVideo(masterId, videoId)
		vh.answerMultipart(writer, videoBytes, http.StatusOK, err)
	}
}

func (vh *VideosHandlers) GetVideoById (writer http.ResponseWriter, req *http.Request) {
	vh.getVideo(writer, req, false)
}

func (vh *VideosHandlers) GetIntro (writer http.ResponseWriter, req *http.Request) {
	vh.getVideo(writer, req, true)
}

func (vh *VideosHandlers) ChangeIntro (writer http.ResponseWriter, req *http.Request) {
	var err error
	var videoData models.VideoData
	sent, masterId := vh.validateMasterId(writer, req)
	if sent {
		return
	}
	sent, file, err := vh.getFileFromFormData(writer, req)
	if sent {
		return
	}
	err = vh.VideosUC.ChangeMasterIntro(&videoData, file, masterId)
	vh.answerIntroPut(writer, videoData, http.StatusOK, err)
}

func (vh *VideosHandlers) deleteVideo (writer http.ResponseWriter, req *http.Request, intro bool) {
	var err error
	sent, masterId := vh.validateMasterId(writer, req)
	if sent {
		return
	}
	if intro {
		err = vh.VideosUC.DeleteMasterIntro(masterId)
	} else {
		sent, videoId := vh.validateVideoId(writer, req)
		if sent {
			return
		}
		err = vh.VideosUC.DeleteMasterVideo(masterId, videoId)
	}
	vh.answerEmpty(writer, http.StatusOK, err)
}

func (vh *VideosHandlers) DeleteVideoById (writer http.ResponseWriter, req *http.Request) {
	vh.deleteVideo(writer, req, false)
}

func (vh *VideosHandlers) DeleteIntro (writer http.ResponseWriter, req *http.Request) {
	vh.deleteVideo(writer, req, true)
}

func (vh *VideosHandlers) GetVideoDataById(writer http.ResponseWriter, req *http.Request) {
	var err error
	sent, masterId := vh.validateMasterId(writer, req)
	if sent {
		return
	}
	sent, videoId := vh.validateVideoId(writer, req)
	if sent {
		return
	}
	videoData := models.VideoData{
		Id: videoId,
	}
	err = vh.VideosUC.GetVideoDataById(&videoData, masterId)
	vh.answerVideo(writer, videoData, http.StatusOK, err)
}

func (vh *VideosHandlers) GetIntroData(writer http.ResponseWriter, req *http.Request) {
	var err error
	sent, masterId := vh.validateMasterId(writer, req)
	if sent {
		return
	}
	videoData := models.VideoData{
		Intro: true,
	}
	err = vh.VideosUC.GetIntroData(&videoData, masterId)
	vh.answerIntroPut(writer, videoData, http.StatusOK, err)
}

func (vh *VideosHandlers) ChangeVideoData(writer http.ResponseWriter, req *http.Request) {
	var err error
	sent, masterId := vh.validateMasterId(writer, req)
	if sent {
		return
	}
	sent, videoId := vh.validateVideoId(writer, req)
	if sent {
		return
	}
	var videoData models.VideoData
	err = json.UnmarshalFromReader(req.Body, &videoData)
	if err != nil {
		jsonError := fmt.Errorf("error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(jsonError.Error()))
		return
	}
	if videoId != videoData.Id {
		paramError := fmt.Errorf("wrong video id parameter")
		logger.Errorf(paramError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(paramError.Error()))
		return
	}
	err = vh.VideosUC.ChangeVideoData(&videoData, masterId)
	vh.answerVideo(writer, videoData, http.StatusOK, err)
}

func (vh *VideosHandlers) ChangeIntroData(writer http.ResponseWriter, req *http.Request) {
	var err error
	sent, masterId := vh.validateMasterId(writer, req)
	if sent {
		return
	}
	var videoData models.VideoData
	err = json.UnmarshalFromReader(req.Body, &videoData)
	if err != nil {
		jsonError := fmt.Errorf("error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(jsonError.Error()))
		return
	}
	videoData.Intro = true
	err = vh.VideosUC.ChangeIntroData(&videoData, masterId)
	vh.answerIntroPut(writer, videoData, http.StatusOK, err)
}

func (vh *VideosHandlers) answerIntroPost(writer http.ResponseWriter, videoData models.VideoData, statusCode int, err error) {
	sent := vh.handlers.handleErrorConflict(writer, err)
	if !sent {
		utils.CreateAnswerVideoDataJson(writer, statusCode, videoData)
	}
}

func (vh *VideosHandlers) answerIntroPut(writer http.ResponseWriter, videoData models.VideoData, statusCode int, err error) {
	sent := vh.handlers.handleErrorNoContent(writer, err)
	if !sent {
		utils.CreateAnswerVideoDataJson(writer, statusCode, videoData)
	}
}

func (vh *VideosHandlers) answerVideo(writer http.ResponseWriter, videoData models.VideoData, statusCode int, err error) {
	sent := vh.handlers.handleError(writer, err)
	if !sent {
		utils.CreateAnswerVideoDataJson(writer, statusCode, videoData)
	}
}

func (vh *VideosHandlers) answerVideos(writer http.ResponseWriter, videoData []models.VideoData, err error) {
	sent := vh.handlers.handleError(writer, err)
	if !sent {
		utils.CreateAnswerVideosDataJson(writer, http.StatusOK, videoData)
	}
}

func (vh *VideosHandlers) answerMultipart(writer http.ResponseWriter, video []byte, statusCode int, err error) {
	sent := vh.handlers.handleError(writer, err)
	if !sent {
		utils.CreateAnswerMultipart(writer, statusCode, video)
	}
}

func (vh *VideosHandlers) answerMultipartIntro(writer http.ResponseWriter, video []byte, statusCode int, err error) {
	sent := vh.handlers.handleErrorNoContent(writer, err)
	if !sent {
		utils.CreateAnswerMultipart(writer, statusCode, video)
	}
}

func (vh *VideosHandlers) answerEmpty(writer http.ResponseWriter, statusCode int, err error) {
	sent := vh.handlers.handleError(writer, err)
	if !sent {
		utils.CreateEmptyBodyAnswer(writer, statusCode)
	}
}