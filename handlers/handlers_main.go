package handlers

import (
	"errors"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"net/http"
)

type Handlers struct {
	UsersHandlers     *UsersHandlers
	MastersHandlers   *MastersHandlers
	StudentsHandlers  *StudentsHandlers
	LanguagesHandlers *LanguagesHandlers
	ThemesHandlers    *ThemesHandlers
	VideosHandlers    *VideosHandlers
	AvatarsHandlers   *AvatarsHandlers
	badRequestError *models.BadRequestError
	conflictError *models.ConflictError
	noContentError *models.NoContentError
}

func (handlers *Handlers) Init(usersUC useCases.UsersUCInterface, mastersUC useCases.MastersUCInterface, studentsUC useCases.StudentsUCInterface,
	themesUC useCases.ThemesUCInterface, languagesUC useCases.LanguagesUCInterface,
	videosUC useCases.VideosUCInterface, avatarsUC useCases.AvatarsUCInterface) error {
	handlers.UsersHandlers = &UsersHandlers{handlers, usersUC}
	handlers.MastersHandlers = &MastersHandlers{handlers, mastersUC}
	handlers.StudentsHandlers = &StudentsHandlers{handlers, studentsUC}
	handlers.LanguagesHandlers = &LanguagesHandlers{handlers, languagesUC}
	handlers.ThemesHandlers = &ThemesHandlers{handlers, themesUC}
	handlers.AvatarsHandlers = &AvatarsHandlers{handlers, avatarsUC}
	handlers.VideosHandlers = &VideosHandlers{
		handlers:handlers,
		VideosUC:videosUC,
	VideoParseConfig:VideoParseConfig{
		FormDataKey:  "video",
		VideoFormats: map[string]bool{
			"video/webm":               true,
			"audio/ogg":                true,
			"video/mp4":                true,
			"video/quicktime":          true,
			"video/x-msvideo":          true,
			"application/octet-stream": true,
		},
	}}
	return nil
}

func (handlers *Handlers) handleErrorConflict(writer http.ResponseWriter, err error) bool {
	if errors.As(err, &handlers.conflictError) {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusConflict, models.CreateMessage(err.Error()))
		return true
	}
	return handlers.handleError(writer, err)
}

func (handlers *Handlers) handleErrorNoContent(writer http.ResponseWriter, err error) bool {
	if errors.As(err, &handlers.noContentError) {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusConflict, models.CreateMessage(err.Error()))
		return true
	}
	return handlers.handleError(writer, err)
}

func (handlers *Handlers) handleError(writer http.ResponseWriter, err error) bool {
	if errors.As(err, &handlers.badRequestError) {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(err.Error()))
		return true
	}
	if err != nil {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(err.Error()))
		return true
	}
	return false
}