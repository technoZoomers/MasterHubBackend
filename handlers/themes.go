package handlers

import (
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"net/http"
)

type ThemesHandlers struct {
	handlers *Handlers
	ThemesUC useCases.ThemesUCInterface
}


func (th *ThemesHandlers) Get(writer http.ResponseWriter, req *http.Request) {
	themes, err := th.ThemesUC.Get()
	th.answerThemes(writer, themes, err)
}

func (th *ThemesHandlers) GetThemeById(writer http.ResponseWriter, req *http.Request) {
	sent, themeId := th.handlers.validateThemeId(writer, req)
	if sent {
		return
	}
	var theme models.Theme
	theme.Id = themeId
	err := th.ThemesUC.GetThemeById(&theme)
	th.answerTheme(writer, theme, err)
}

func (th *ThemesHandlers) answerTheme(writer http.ResponseWriter, theme models.Theme, err error) {
	sent := th.handlers.handleError(writer, err)
	if !sent {
		utils.CreateAnswerThemeJson(writer, http.StatusOK, theme)
	}
}

func (th *ThemesHandlers) answerThemes(writer http.ResponseWriter, themes []models.Theme, err error) {
	sent := th.handlers.handleError(writer, err)
	if !sent {
		utils.CreateAnswerThemesJson(writer, http.StatusOK, themes)
	}
}
