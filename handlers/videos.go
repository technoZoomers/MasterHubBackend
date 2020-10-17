package handlers

import (
	"errors"
	"fmt"
	"github.com/google/logger"
	"github.com/gorilla/mux"
	json "github.com/mailru/easyjson"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"net/http"
	"strconv"
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

func (vh *VideosHandlers) Upload(writer http.ResponseWriter, req *http.Request) {
	vh.uploadVideo(writer, req, false)
}

func (vh *VideosHandlers) UploadIntro(writer http.ResponseWriter, req *http.Request) {
	vh.uploadVideo(writer, req, true)
}

func (vh *VideosHandlers) uploadVideo (writer http.ResponseWriter, req *http.Request, intro bool) {
	var err error
	var videoData models.VideoData
	masterIdString := mux.Vars(req)["id"]
	masterId, err := strconv.ParseInt(masterIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting master id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}
	file, fileHeader, err := req.FormFile(vh.VideoParseConfig.FormDataKey)
	if err != nil {
		parseError := fmt.Errorf("error parsing video: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(parseError.Error()))
		return
	}
	if !vh.VideoParseConfig.VideoFormats[fileHeader.Header.Get("Content-Type")] {
		parseError := fmt.Errorf("wrong mime type:%s, expected video", fileHeader.Header.Get("Content-Type"))
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(parseError.Error()))
		return
	}
	if intro {
		err = vh.VideosUC.NewMasterIntro(&videoData, file, masterId)
		vh.answerIntro(writer, videoData, http.StatusCreated, err)
	} else {
		err = vh.VideosUC.NewMasterVideo(&videoData, file, masterId)
		vh.answerVideo(writer, videoData, http.StatusCreated, err)
	}
}

func (vh *VideosHandlers) GetVideosByMasterId(writer http.ResponseWriter, req *http.Request) {
	var err error
	masterIdString := mux.Vars(req)["id"]
	masterId, err := strconv.ParseInt(masterIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting master id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}
	videos, err := vh.VideosUC.GetVideosByMasterId(masterId)
	vh.answerVideos(writer, videos, err)
}

func (vh *VideosHandlers) getVideo(writer http.ResponseWriter, req *http.Request, intro bool) {
	var err error
	masterIdString := mux.Vars(req)["id"]
	masterId, err := strconv.ParseInt(masterIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting master id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}
	videoIdString := mux.Vars(req)["videoId"]
	videoId, err := strconv.ParseInt(videoIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting video id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}

	videoBytes, err := vh.VideosUC.GetMasterVideo(masterId, videoId)

	if errors.As(err, &vh.handlers.badRequestError) {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(err.Error()))
		return
	}
	if err != nil {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(err.Error()))
		return
	}

	writer.Header().Set("content-type", "multipart/form-data;boundary=1")
	_, err = writer.Write(videoBytes)
	if err != nil {
		writeBytesError := fmt.Errorf("error writing video bytes to response: %v", err.Error())
		logger.Error(writeBytesError)
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(writeBytesError.Error()))
		return
	}
	utils.CreateEmptyBodyAnswerJson(writer, http.StatusOK)
}

func (vh *VideosHandlers) GetVideoById(writer http.ResponseWriter, req *http.Request) {
	vh.getVideo(writer, req, false)
}

func (vh *VideosHandlers) GetIntro(writer http.ResponseWriter, req *http.Request) {
	vh.getVideo(writer, req, true)
}

func (vh *VideosHandlers) GetVideoDataById(writer http.ResponseWriter, req *http.Request) {
	var err error
	masterIdString := mux.Vars(req)["id"]
	masterId, err := strconv.ParseInt(masterIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting master id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}
	videoIdString := mux.Vars(req)["videoId"]
	videoId, err := strconv.ParseInt(videoIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting video id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}
	videoData := models.VideoData{
		Id: videoId,
	}
	err = vh.VideosUC.GetVideoDataById(&videoData, masterId)
	vh.answerVideo(writer, videoData, http.StatusOK, err)
}

func (vh *VideosHandlers) ChangeVideoData(writer http.ResponseWriter, req *http.Request) {
	var err error
	masterIdString := mux.Vars(req)["id"]
	masterId, err := strconv.ParseInt(masterIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting master id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}
	videoIdString := mux.Vars(req)["videoId"]
	videoId, err := strconv.ParseInt(videoIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting video id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
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

func (vh *VideosHandlers) answerIntro(writer http.ResponseWriter, videoData models.VideoData, statusCode int, err error) {
	sent := vh.handlers.handleErrorConflict(writer, err)
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