package handlers

import (
	"fmt"
	"github.com/google/logger"
	json "github.com/mailru/easyjson"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"net/http"
)

type UsersHandlers struct {
	handlers *Handlers
	UsersUC  useCases.UsersUCInterface
}

func (uh *UsersHandlers) GetUserById(writer http.ResponseWriter, req *http.Request) {
	sent, userId := uh.handlers.validateUserId(writer, req)
	if sent {
		return
	}
	var user models.User
	user.Id = userId
	err := uh.UsersUC.GetUserById(&user)
	uh.answerUser(writer, user, err)
}

func (uh *UsersHandlers) Login(writer http.ResponseWriter, req *http.Request) {
	var user models.User
	err := json.UnmarshalFromReader(req.Body, &user)
	if err != nil {
		jsonError := fmt.Errorf("error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(jsonError.Error()))
		return
	}
	err = uh.UsersUC.Login(&user)
	uh.answerUser(writer, user, err)
}

func (uh *UsersHandlers) answerUser(writer http.ResponseWriter, user models.User, err error) {
	sent := uh.handlers.handleError(writer, err)
	if !sent {
		utils.CreateAnswerUserJson(writer, http.StatusOK, user)
	}
}