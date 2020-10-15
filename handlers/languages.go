package handlers

import (
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"net/http"
)

type LanguagesHandlers struct {
	LanguagesUC useCases.LanguagesUCInterface
}

func (lh *LanguagesHandlers) Get(writer http.ResponseWriter, req *http.Request) {
	languages, err := lh.LanguagesUC.Get()
	if err != nil {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(err.Error()))
		return
	}
	utils.CreateAnswerLanguagesJson(writer, http.StatusOK, languages)
}
