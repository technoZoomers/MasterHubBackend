package handlers

import (
	"fmt"
	"github.com/google/logger"
	"github.com/google/uuid"
	json "github.com/mailru/easyjson"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"net/http"
	"time"
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
	sent = uh.handlers.checkUserAuth(writer, req, userId)
	if sent {
		return
	}
	var user models.User
	user.Id = userId
	err := uh.UsersUC.GetUserById(&user)
	uh.answerUser(writer, user, err)
}

func (uh *UsersHandlers) Login(writer http.ResponseWriter, req *http.Request) {
	sent := uh.handlers.checkNoAuth(writer, req)
	if sent {
		return
	}
	var user models.User
	err := json.UnmarshalFromReader(req.Body, &user)
	if err != nil {
		jsonError := fmt.Errorf("error unmarshaling json: %v", err.Error())
		logger.Errorf(jsonError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(jsonError.Error()))
		return
	}
	err = uh.UsersUC.Login(&user)
	var cookie http.Cookie
	if err == nil {
		err = uh.setCookie(&user, &cookie)
		if err != nil {
			cookieError := fmt.Errorf("error setting cookie: %v", err.Error())
			logger.Errorf(cookieError.Error())
			err = cookieError
		}
	}
	uh.answerUserLogin(writer, user, &cookie, err)
}

func (uh *UsersHandlers) Logout(writer http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie(uh.handlers.cookieString)
	if err != nil {
		cookieError := fmt.Errorf("error finding cookie: %v", err.Error())
		logger.Errorf(cookieError.Error())
		uh.answerEmptyLogout(writer, cookie, cookieError)
	} else {
		uh.answerEmptyLogout(writer, cookie, uh.deleteCookie(cookie))
	}
}

func (uh *UsersHandlers) setCookie(user *models.User, cookie *http.Cookie) error {
	token := uuid.New()
	cookie.Name = uh.handlers.cookieString
	cookie.Value = token.String()
	cookie.Expires = time.Now().Add(365 * 24 * time.Hour)
	cookie.SameSite = http.SameSiteNoneMode
	cookie.Secure = true
	//cookie.HttpOnly = true
	cookie.Path = "/"
	return uh.UsersUC.InsertCookie(user.Id, cookie.Value)
}


func (uh *UsersHandlers) deleteCookie(cookie *http.Cookie) error {
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	return uh.UsersUC.DeleteCookie(cookie.Value)
}

func (uh *UsersHandlers) CheckAuth(writer http.ResponseWriter, req *http.Request) {
	user, ok := req.Context().Value(uh.handlers.contextUserKey).(models.User)
	if !ok {
		internalError := fmt.Errorf("error getting value from context")
		logger.Errorf(internalError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(internalError.Error()))
		return
	}
	uh.answerUser(writer, user, nil)
}

func (uh *UsersHandlers) answerUser(writer http.ResponseWriter, user models.User, err error) {
	sent := uh.handlers.handleError(writer, err)
	if !sent {
		utils.CreateAnswerUserJson(writer, http.StatusOK, user)
	}
}

func (uh *UsersHandlers) answerUserLogin(writer http.ResponseWriter, user models.User, cookie *http.Cookie, err error) {
	sent := uh.handlers.handleError(writer, err)
	if !sent {
		http.SetCookie(writer, cookie)
		utils.CreateAnswerUserJson(writer, http.StatusOK, user)
	}
}

func (uh *UsersHandlers) answerEmptyLogout(writer http.ResponseWriter, cookie *http.Cookie, err error) {
	sent := uh.handlers.handleError(writer, err)
	if !sent {
		http.SetCookie(writer, cookie)
		utils.CreateEmptyBodyAnswer(writer, http.StatusOK)
	}
}