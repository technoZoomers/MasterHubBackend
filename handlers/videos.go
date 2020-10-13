package handlers

import (
	"bufio"
	"fmt"
	"github.com/google/logger"
	"github.com/gorilla/mux"
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
	err = req.ParseMultipartForm(32 << 20)
	if err != nil {
		parseError := fmt.Errorf("error parsing video: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(parseError.Error()))
		return
	}
	file, fileHeader, err := req.FormFile(utils.FORM_DATA_VIDEO_KEY)
	if err != nil {
		parseError := fmt.Errorf("error parsing video: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(parseError.Error()))
		return
	}
	if fileHeader.Header.Get("Content-Type") != utils.VIDEO_FORMAT {
		parseError := fmt.Errorf("wrong mime type, expected video")
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

	userNotExist, file, err := ph.PicUC.GetUserPic(&user)
	if userNotExist {
		network.CreateErrorAnswerJson(writer, utils.StatusCode("Bad Request"), models.CreateMessage(err.Error()))
		return
	}
	if err != nil {
		logger.Error(err)
		network.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(err.Error()))
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()

	bytes := make([]byte, size)
	_, err = reader.Read(bytes)

	writer.Header().Set("content-type", "multipart/form-data;boundary=1")

	_, err = writer.Write(bytes)
	if err != nil {
		logger.Error(err)
		network.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(err.Error()))
		return
	}
}