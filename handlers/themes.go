package handlers

import (
	"errors"
	"fmt"
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"net/http"
	"strconv"
)

type ThemesHandlers struct {
	handlers     *Handlers
	ThemesUC useCases.ThemesUCInterface
}

func (th *ThemesHandlers) Get(writer http.ResponseWriter, req *http.Request) {
	themes, err := th.ThemesUC.Get()
	if err != nil {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(err.Error()))
		return
	}
	utils.CreateAnswerThemesJson(writer, http.StatusOK, themes)
}

func (th *ThemesHandlers) GetThemeById(writer http.ResponseWriter, req *http.Request) {
	themeIdString := mux.Vars(req)["id"]
	themeId, err := strconv.ParseInt(themeIdString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting theme id parameter: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return
	}
	var theme models.Theme
	theme.Id = themeId
	err = th.ThemesUC.GetThemeById(&theme)
	if errors.As(err, &th.handlers.badRequestError) {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(err.Error()))
		return
	}
	if err != nil {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(err.Error()))
		return
	}
	utils.CreateAnswerThemeJson(writer, http.StatusOK, theme)
}
