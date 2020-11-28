package handlers

import (
	"errors"
	"fmt"
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v2"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"net/http"
	"strconv"
)

type Handlers struct {
	UsersHandlers          *UsersHandlers
	MastersHandlers        *MastersHandlers
	StudentsHandlers       *StudentsHandlers
	LanguagesHandlers      *LanguagesHandlers
	ThemesHandlers         *ThemesHandlers
	VideosHandlers         *VideosHandlers
	AvatarsHandlers        *AvatarsHandlers
	ChatsHandlers          *ChatsHandlers
	WSHandlers             *WSHandlers
	LessonsHandlers        *LessonsHandlers
	VCHandlers             *VCHandlers
	badRequestError        *models.BadRequestError
	conflictError          *models.ConflictError
	noContentError         *models.NoContentError
	badQueryParameterError *models.BadQueryParameterError
	forbiddenError         *models.ForbiddenError
	notAcceptableError     *models.NotAcceptableError
	contextUserKey         string
	contextCookieKey       string
	cookieString           string
	contextAuthorisedKey   string
}

func (handlers *Handlers) Init(usersUC useCases.UsersUCInterface, mastersUC useCases.MastersUCInterface, studentsUC useCases.StudentsUCInterface,
	themesUC useCases.ThemesUCInterface, languagesUC useCases.LanguagesUCInterface,
	videosUC useCases.VideosUCInterface, avatarsUC useCases.AvatarsUCInterface,
	chatsUC useCases.ChatsUCInterface, wsUC useCases.WebsocketsUCInterface, lessonsUC useCases.LessonsUCInterface,
	videocallsUC useCases.VideocallsUCInterface) error {
	handlers.UsersHandlers = &UsersHandlers{
		handlers: handlers,
		UsersUC:  usersUC,
		wsUC:     wsUC,
		vcUC:     videocallsUC,
	}
	handlers.MastersHandlers = &MastersHandlers{
		handlers:  handlers,
		MastersUC: mastersUC,
		MastersQueryKeys: MastersQueryKeys{
			Subtheme:        "subtheme",
			Theme:           "theme",
			Qualification:   "qualification",
			EducationFormat: "educationFormat",
			Language:        "language",
			Search:          "search",
			Limit:           "limit",
			Offset:          "offset",
		},
	}
	handlers.StudentsHandlers = &StudentsHandlers{handlers, studentsUC}
	handlers.LanguagesHandlers = &LanguagesHandlers{handlers, languagesUC}
	handlers.ThemesHandlers = &ThemesHandlers{handlers, themesUC}
	handlers.AvatarsHandlers = &AvatarsHandlers{
		handlers:  handlers,
		AvatarsUC: avatarsUC,
		AvatarParseConfig: AvatarParseConfig{
			FormDataKey: "avatar",
			ImgFormats: map[string]bool{
				"image/jpeg":    true,
				"image/png":     true,
				"image/gif":     true,
				"image/svg+xml": true,
			},
		},
	}
	handlers.VideosHandlers = &VideosHandlers{
		handlers: handlers,
		VideosUC: videosUC,
		VideoParseConfig: VideoParseConfig{
			FormDataKey: "video",
			VideoFormats: map[string]bool{
				"video/webm":               true,
				"audio/ogg":                true,
				"video/mp4":                true,
				"video/quicktime":          true,
				"video/x-msvideo":          true,
				"application/octet-stream": true,
			},
		},
		VideosQueryKeys: VideosQueryKeys{
			Subtheme: "subtheme",
			Theme:    "theme",
			Old:      "old",
			Popular:  "popular",
			Limit:    "limit",
			Offset:   "offset",
		}}
	handlers.ChatsHandlers = &ChatsHandlers{
		handlers: handlers,
		ChatsUC:  chatsUC,
		ChatsQueryKeys: ChatsQueryKeys{
			Type:   "type",
			Limit:  "limit",
			Offset: "offset",
		},
	}
	handlers.WSHandlers = &WSHandlers{
		handlers: handlers,
		ChatsUC:  chatsUC,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool { // TODO: fix this
				return true
			},
		},
		WebsocketsUC: wsUC,
	}
	handlers.LessonsHandlers = &LessonsHandlers{
		handlers:  handlers,
		LessonsUC: lessonsUC,
		LessonsQueryKeys: LessonsQueryKeys{
			Status: "status",
		},
	}
	mediaEngine := webrtc.MediaEngine{}
	mediaEngine.RegisterCodec(webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeVP8, 90000))
	mediaEngine.RegisterCodec(webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000))
	handlers.VCHandlers = &VCHandlers{
		handlers:     handlers,
		videocallsUC: videocallsUC,
		wsUC:         wsUC,
		webrtcConfig: webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs: []string{"stun:stun.l.google.com:19302"},
				},
			},
		},
		webrtcAPI: webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine)),
	}
	handlers.cookieString = "user_session"
	handlers.contextCookieKey = "cookie_key"
	handlers.contextUserKey = "user_key"
	handlers.contextAuthorisedKey = "auth_key"
	return nil
}

func (handlers *Handlers) handleErrorBadQueryParameter(writer http.ResponseWriter, err error) bool {
	if errors.As(err, &handlers.badQueryParameterError) {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusBadRequest, models.CreateMessage(err.Error()))
		return true
	}
	return handlers.handleError(writer, err)
}

func (handlers *Handlers) handleErrorConflict(writer http.ResponseWriter, err error) bool {
	if errors.As(err, &handlers.conflictError) {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusConflict, models.CreateMessage(err.Error()))
		return true
	}
	return handlers.handleForbiddenError(writer, err)
}

func (handlers *Handlers) handleErrorNoContent(writer http.ResponseWriter, err error) bool {
	if errors.As(err, &handlers.noContentError) {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusNoContent, models.CreateMessage(err.Error()))
		return true
	}
	return handlers.handleError(writer, err)
}

func (handlers *Handlers) handleForbiddenError(writer http.ResponseWriter, err error) bool {
	if errors.As(err, &handlers.forbiddenError) {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusForbidden, models.CreateMessage(err.Error()))
		return true
	}
	return handlers.handleError(writer, err)
}

func (handlers *Handlers) handleNotAcceptableError(writer http.ResponseWriter, err error) bool {
	if errors.As(err, &handlers.notAcceptableError) {
		logger.Error(err)
		utils.CreateErrorAnswerJson(writer, http.StatusNotAcceptable, models.CreateMessage(err.Error()))
		return true
	}
	return handlers.handleErrorConflict(writer, err)
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

func (handlers *Handlers) validateUserId(writer http.ResponseWriter, req *http.Request) (bool, int64) {
	return handlers.validateId(writer, req, "id", "user")
}

func (handlers *Handlers) validatePeerId(writer http.ResponseWriter, req *http.Request) (bool, int64) {
	return handlers.validateId(writer, req, "peerId", "peer")
}

func (handlers *Handlers) validateStudentId(writer http.ResponseWriter, req *http.Request) (bool, int64) {
	return handlers.validateId(writer, req, "id", "student")
}

func (handlers *Handlers) validateMasterId(writer http.ResponseWriter, req *http.Request) (bool, int64) {
	return handlers.validateId(writer, req, "id", "master")
}

func (handlers *Handlers) validateChatId(writer http.ResponseWriter, req *http.Request) (bool, int64) {
	return handlers.validateId(writer, req, "chatId", "chat")
}
func (handlers *Handlers) validateChatIdSimple(writer http.ResponseWriter, req *http.Request) (bool, int64) {
	return handlers.validateId(writer, req, "id", "chat")
}

func (handlers *Handlers) validateLessonId(writer http.ResponseWriter, req *http.Request) (bool, int64) {
	return handlers.validateId(writer, req, "lessonId", "lesson")
}

func (handlers *Handlers) validateThemeId(writer http.ResponseWriter, req *http.Request) (bool, int64) {
	return handlers.validateId(writer, req, "id", "theme")
}

func (handlers *Handlers) validateLessonIdSimple(writer http.ResponseWriter, req *http.Request) (bool, int64) {
	return handlers.validateId(writer, req, "id", "lesson")
}

func (handlers *Handlers) validateVideoId(writer http.ResponseWriter, req *http.Request) (bool, int64) {
	return handlers.validateId(writer, req, "videoId", "video")
}

func (handlers *Handlers) checkNoAuth(writer http.ResponseWriter, r *http.Request) bool {
	auth, ok := r.Context().Value(handlers.contextAuthorisedKey).(bool)
	if !ok {
		internalError := fmt.Errorf("error getting value from context")
		logger.Errorf(internalError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(internalError.Error()))
		return true
	}
	if auth {
		authError := fmt.Errorf("user already logged in")
		logger.Errorf(authError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusForbidden, models.CreateMessage(authError.Error()))
		return true
	}
	return false
}

func (handlers *Handlers) checkUserAuth(writer http.ResponseWriter, r *http.Request, userId int64) bool {
	user, ok := r.Context().Value(handlers.contextUserKey).(models.User)
	if !ok {
		internalError := fmt.Errorf("error getting value from context")
		logger.Errorf(internalError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(internalError.Error()))
		return true
	}
	if user.Id != userId {
		authError := fmt.Errorf("can't get another users info")
		logger.Errorf(authError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusForbidden, models.CreateMessage(authError.Error()))
		return true
	}
	return false
}

func (handlers *Handlers) checkChatAuth(writer http.ResponseWriter, r *http.Request, chatId int64) bool {
	user, ok := r.Context().Value(handlers.contextUserKey).(models.User)
	if !ok {
		internalError := fmt.Errorf("error getting value from context")
		logger.Errorf(internalError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(internalError.Error()))
		return true
	}
	err := handlers.ChatsHandlers.ChatsUC.CheckChatByUserId(chatId, user.Id)
	if err != nil {
		if errors.As(err, &handlers.badRequestError) {
			authError := fmt.Errorf("can't get chat with user id and chat id")
			logger.Errorf(authError.Error())
			utils.CreateErrorAnswerJson(writer, http.StatusForbidden, models.CreateMessage(authError.Error()))
			return true
		} else {
			logger.Errorf(err.Error())
			utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(err.Error()))
			return true
		}
	}
	return false
}
