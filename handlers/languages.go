package handlers

import (
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"net/http"
)

type LanguagesHandlers struct {
	LanguagesUC useCases.LanguagesUCInterface
}

func (lh *LanguagesHandlers) Get(writer http.ResponseWriter, req *http.Request) {
	languages, err := lh.LanguagesUC.GetLanguages()
	if err != nil {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, utils.StatusCode("Internal Server Error"), models.CreateMessage(err.Error()))
		return
	}
	utils.CreateAnswerTransactionsJson(writer, utils.StatusCode("OK"), txs)
}