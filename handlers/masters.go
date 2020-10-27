package handlers

import (
	"fmt"
	"github.com/google/logger"
	json "github.com/mailru/easyjson"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"net/http"
	"net/url"
	"strconv"
)

type MastersHandlers struct {
	handlers         *Handlers
	MastersUC        useCases.MastersUCInterface
	MastersQueryKeys MastersQueryKeys
}

type MastersQueryKeys struct {
	Subtheme        string
	Theme           string
	Qualification   string
	EducationFormat string
	Language        string
	Search          string
	Limit           string
	Offset          string
}

func (mh *MastersHandlers) GetMasterById(writer http.ResponseWriter, req *http.Request) {
	sent, masterId := mh.handlers.validateMasterId(writer, req)
	if sent {
		return
	}
	var master models.Master
	master.UserId = masterId
	err := mh.MastersUC.GetMasterById(&master)
	mh.answerMaster(writer, master, err)
}

func (mh *MastersHandlers) ChangeMasterData(writer http.ResponseWriter, req *http.Request) {
	sent, masterId := mh.handlers.validateMasterId(writer, req)
	if sent {
		return
	}
	sent = mh.handlers.checkUserAuth(writer, req, masterId)
	if sent {
		return
	}
	var master models.Master
	err := json.UnmarshalFromReader(req.Body, &master)
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
	mh.answerMaster(writer, master, err)
}

func (mh *MastersHandlers) parseMastersQuery(query url.Values, mastersQuery *models.MastersQueryValues) error {
	mastersQuery.Subtheme = query[mh.MastersQueryKeys.Subtheme]
	mastersQuery.Theme = query.Get(mh.MastersQueryKeys.Theme)
	mastersQuery.EducationFormat = query.Get(mh.MastersQueryKeys.EducationFormat)
	mastersQuery.Qualification = query.Get(mh.MastersQueryKeys.Qualification)
	mastersQuery.Language = query[mh.MastersQueryKeys.Language]
	mastersQuery.Search = query.Get(mh.MastersQueryKeys.Search)
	limitString := query.Get(mh.MastersQueryKeys.Limit)
	if limitString != "" {
		limit, err := strconv.ParseInt(limitString, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing query parameter %s: %v", mh.MastersQueryKeys.Limit, err.Error())
		}
		mastersQuery.Limit = limit
	}
	offsetString := query.Get(mh.MastersQueryKeys.Offset)
	if offsetString != "" {
		offset, err := strconv.ParseInt(offsetString, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing query parameter %s: %v", mh.MastersQueryKeys.Offset, err.Error())
		}
		mastersQuery.Offset = offset
	}
	return nil
}

func (mh *MastersHandlers) Get(writer http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	var mastersQuery models.MastersQueryValues
	err := mh.parseMastersQuery(query, &mastersQuery)
	if err != nil {
		logger.Errorf(err.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(err.Error()))
		return
	}
	masters, err := mh.MastersUC.Get(mastersQuery)
	mh.answerMasters(writer, masters, err)
}

func (mh *MastersHandlers) Register(writer http.ResponseWriter, req *http.Request) {
	var newMaster models.MasterFull
	err := json.UnmarshalFromReader(req.Body, &newMaster)
	if err != nil {
		jsonError := fmt.Errorf("error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(jsonError.Error()))
		return
	}
	err = mh.MastersUC.Register(&newMaster)
	mh.answerMasterFull(writer, newMaster, err)

}

func (mh *MastersHandlers) answerMaster(writer http.ResponseWriter, master models.Master, err error) {
	sent := mh.handlers.handleErrorConflict(writer, err)
	if !sent {
		utils.CreateAnswerMasterJson(writer, http.StatusOK, master)
	}
}

func (mh *MastersHandlers) answerMasterFull(writer http.ResponseWriter, masterFull models.MasterFull,  err error) {
	sent := mh.handlers.handleErrorConflict(writer, err)
	if !sent {
		utils.CreateAnswerMasterFullJson(writer, http.StatusCreated, masterFull)
	}
}

func (mh *MastersHandlers) answerMasters(writer http.ResponseWriter, masters models.Masters, err error) {
	sent := mh.handlers.handleErrorBadQueryParameter(writer, err)
	if !sent {
		utils.CreateAnswerMastersJson(writer, http.StatusOK, masters)
	}
}
