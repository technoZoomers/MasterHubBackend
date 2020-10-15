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

type MastersHandlers struct {
	MastersUC useCases.MastersUCInterface
}

func (mh *MastersHandlers) GetMasterById(writer http.ResponseWriter, req *http.Request) {
	masterIdString := mux.Vars(req)["id"]
	masterId, err := strconv.ParseInt(masterIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting master id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}
	var master models.Master
	master.UserId = masterId
	err = mh.MastersUC.GetMasterById(&master)
	var badReqError *models.BadRequestError
	if errors.As(err, &badReqError) {
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

func (mh *MastersHandlers) ChangeMasterData(writer http.ResponseWriter, req *http.Request) {
	masterIdString := mux.Vars(req)["id"]
	masterId, err := strconv.ParseInt(masterIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting master id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}
	var master models.Master
	err = json.UnmarshalFromReader(req.Body, &master)
	if err != nil {
		jsonError := fmt.Errorf("error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(jsonError.Error()))
		return
	}
	if masterId != master.UserId {
		paramError := fmt.Errorf("wrong master id parameter")
		logger.Errorf(paramError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(paramError.Error()))
		return
	}
	err = mh.MastersUC.ChangeMasterData(&master)
	var badReqError *models.BadRequestError
	if errors.As(err, &badReqError) {
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
