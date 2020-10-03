package handlers

import (
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
	file, _, err := req.FormFile("video")
	if err != nil {
		parseError := fmt.Errorf("error parsing video: %v", err.Error())
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