package handlers

import (
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
	VideosUC useCases.VideosUCInterface
}

func (vh *VideosHandlers) Upload(writer http.ResponseWriter, req *http.Request) {
	var err error
	var videoData models.VideoData
	masterIdString := mux.Vars(req)["id"]
	masterId, err :=  strconv.ParseInt(masterIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting master id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}
	//err = req.ParseMultipartForm(32 << 20)
	//if err != nil {
	//	parseError := fmt.Errorf("error parsing video: %v", err.Error())
	//	logger.Errorf(parseError.Error())
	//	utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(parseError.Error()))
	//	return
	//}
	file, fileHeader, err := req.FormFile(utils.FORM_DATA_VIDEO_KEY)
	if err != nil {
		parseError := fmt.Errorf("error parsing video: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(parseError.Error()))
		return
	}
	if fileHeader.Header.Get("Content-Type") != utils.VIDEO_FORMAT {
		parseError := fmt.Errorf("wrong mime type:%s, expected video", fileHeader.Header.Get("Content-Type"))
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(parseError.Error()))
		return
	}
	defer file.Close()

	absent, err := vh.VideosUC.NewMasterVideo(&videoData, file, masterId)
	if absent {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(err.Error()))
		return
	}
	if err != nil {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(err.Error()))
		return
	}
	utils.CreateAnswerVideoDataJson(writer, http.StatusOK, videoData)
}

func (vh *VideosHandlers) GetVideosByMasterId(writer http.ResponseWriter, req *http.Request) {
	var err error
	masterIdString := mux.Vars(req)["id"]
	masterId, err :=  strconv.ParseInt(masterIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting master id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}
	videos, absent, err := vh.VideosUC.GetVideosByMasterId(masterId)
	if absent {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(err.Error()))
		return
	}
	if err != nil {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(err.Error()))
		return
	}
	utils.CreateAnswerVideosDataJson(writer, http.StatusOK, videos)
}

func (vh *VideosHandlers) GetVideoById(writer http.ResponseWriter, req *http.Request) {
	var err error
	masterIdString := mux.Vars(req)["id"]
	masterId, err :=  strconv.ParseInt(masterIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting master id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}
	videoIdString := mux.Vars(req)["videoId"]
	videoId, err :=  strconv.ParseInt(videoIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting video id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}

	videoBytes, absent, err := vh.VideosUC.GetMasterVideo(masterId, videoId)
	if absent {
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

func (vh *VideosHandlers) GetVideoDataById(writer http.ResponseWriter, req *http.Request) {
	var err error
	masterIdString := mux.Vars(req)["id"]
	masterId, err :=  strconv.ParseInt(masterIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting master id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}
	videoIdString := mux.Vars(req)["videoId"]
	videoId, err :=  strconv.ParseInt(videoIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting video id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}
	videoData := models.VideoData{
		Id:videoId,
	}
	absent, err := vh.VideosUC.GetVideoDataById(&videoData, masterId)
	if absent {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(err.Error()))
		return
	}
	if err != nil {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(err.Error()))
		return
	}
	utils.CreateAnswerVideoDataJson(writer, http.StatusOK, videoData)
}

func (vh *VideosHandlers) ChangeVideoData (writer http.ResponseWriter, req *http.Request) {
	var err error
	masterIdString := mux.Vars(req)["id"]
	masterId, err :=  strconv.ParseInt(masterIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting master id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}
	videoIdString := mux.Vars(req)["videoId"]
	videoId, err :=  strconv.ParseInt(videoIdString, 10, 64)
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
	absent, err := vh.VideosUC.ChangeVideoData(&videoData, masterId)
	if absent {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(err.Error()))
		return
	}
	if err != nil {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(err.Error()))
		return
	}
	utils.CreateAnswerVideoDataJson(writer, http.StatusOK, videoData)
}