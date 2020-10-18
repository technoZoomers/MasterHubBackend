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
	handlers.MastersHandlers = &MastersHandlers{
		handlers: handlers,
		MastersUC: mastersUC,
		MastersQueryKeys:MastersQueryKeys{
			Subtheme:        "subtheme",
			Theme:           "theme",
			Qualification:   "qualification",
			EducationFormat: "educationFormat",
			Language:        "language",
			Search: "search",
			Limit: "limit",
			Offset: "offset",
		},
	}
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
	if _, ok := err.(*models.ConflictError); ok {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusConflict, models.CreateMessage(err.Error()))
		return true
	}
	//if errors.As(err, &handlers.conflictError) {
	//	logger.Error(err)
	//	utils.CreateErrorAnswerJson(writer, http.StatusConflict, models.CreateMessage(err.Error()))
	//	return true
	//}
	return handlers.handleError(writer, err)
}

func (handlers *Handlers) handleErrorNoContent(writer http.ResponseWriter, err error) bool {
	if _, ok := err.(*models.NoContentError); ok {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusNoContent, models.CreateMessage(err.Error()))
		return true
	}
	//if errors.As(err, &handlers.noContentError) {
	//	logger.Error(err)
	//	utils.CreateErrorAnswerJson(writer, http.StatusNoContent, models.CreateMessage(err.Error()))
	//	return true
	//}
	return handlers.handleError(writer, err)
}

func (handlers *Handlers) handleError(writer http.ResponseWriter, err error) bool {
	if _, ok := err.(*models.BadRequestError); ok {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(err.Error()))
		return true
	}
	//if errors.As(err, &handlers.badRequestError) {
	//	logger.Error(err)
	//	utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(err.Error()))
	//	return true
	//}
	if err != nil {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(err.Error()))
		return true
	}
	return false
}
func (handlers *Handlers) validateId(writer http.ResponseWriter, req *http.Request, idName string, entityName string) (bool, int64) {
	idString := mux.Vars(req)[idName]
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		parseError := fmt.Errorf("error getting %s id parameter: %v", entityName, err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(parseError.Error()))
		return true, id
	}
	return false, id
}