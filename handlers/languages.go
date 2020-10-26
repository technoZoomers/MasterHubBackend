package handlers

import (
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"net/http"
)

type LanguagesHandlers struct {
	handlers    *Handlers
	LanguagesUC useCases.LanguagesUCInterface
}

func (lh *LanguagesHandlers) Get(writer http.ResponseWriter, req *http.Request) {
	languages, err := lh.LanguagesUC.Get()
	lh.answerLanguages(writer, languages, err)
}

func (lh *LanguagesHandlers) answerLanguages(writer http.ResponseWriter, languages models.Languages, err error) {
	sent := lh.handlers.handleError(writer, err)
	if !sent {
		utils.CreateAnswerLanguagesJson(writer, http.StatusOK, languages)
	}
}
