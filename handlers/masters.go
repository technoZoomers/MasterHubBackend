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

type MastersHandlers struct {
	MastersUC useCases.MastersUCInterface
}

func (mh *MastersHandlers) GetMasterById(writer http.ResponseWriter, req *http.Request) {
	masterIdString := mux.Vars(req)["id"]
	masterId, err :=  strconv.ParseInt(masterIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting master id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}
	var master models.Master
	master.UserId = masterId
	absent, err := mh.MastersUC.GetMasterById(&master)
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
	utils.CreateAnswerMasterJson(writer, http.StatusOK, master)
}