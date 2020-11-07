package handlers

import (
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"mime/multipart"
	"net/http"
)

type AvatarsHandlers struct {
	handlers          *Handlers
	AvatarsUC         useCases.AvatarsUCInterface
	AvatarParseConfig AvatarParseConfig
}

type AvatarParseConfig struct {
	FormDataKey string
	ImgFormats  map[string]bool
}

func (ah *AvatarsHandlers) getFileFromFormData(writer http.ResponseWriter, req *http.Request) (bool, multipart.File, error) {
	file, fileHeader, err := req.FormFile(ah.AvatarParseConfig.FormDataKey)
	if err != nil {
		parseError := fmt.Errorf("error parsing avatar: %v", err.Error())
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(parseError.Error()))
		return true, file, err
	}
	if !ah.AvatarParseConfig.ImgFormats[fileHeader.Header.Get("Content-Type")] {
		parseError := fmt.Errorf("wrong mime type:%s, expected image", fileHeader.Header.Get("Content-Type"))
		logger.Errorf(parseError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusInternalServerError, models.CreateMessage(parseError.Error()))
		return true, file, err
	}
	return false, file, err
}

func (ah *AvatarsHandlers) UploadAvatar(writer http.ResponseWriter, req *http.Request) {
	var err error
	sent, userId := ah.handlers.validateUserId(writer, req)
	if sent {
		return
	}
	sent = ah.handlers.checkUserAuth(writer, req, userId)
	if sent {
		return
	}
	sent, file, err := ah.getFileFromFormData(writer, req)
	if sent {
		return
	}
	err = ah.AvatarsUC.NewUserAvatar(file, userId)
	ah.answerEmpty(writer, http.StatusCreated, err)
}

func (ah *AvatarsHandlers) ChangeAvatar(writer http.ResponseWriter, req *http.Request) {
	var err error
	sent, userId := ah.handlers.validateUserId(writer, req)
	if sent {
		return
	}
	sent = ah.handlers.checkUserAuth(writer, req, userId)
	if sent {
		return
	}
	sent, file, err := ah.getFileFromFormData(writer, req)
	if sent {
		return
	}
	err = ah.AvatarsUC.ChangeUserAvatar(file, userId)
	ah.answerEmpty(writer, http.StatusOK, err)
}

func (ah *AvatarsHandlers) GetAvatar(writer http.ResponseWriter, req *http.Request) {
	var err error
	sent, userId := ah.handlers.validateUserId(writer, req)
	if sent {
		return
	}
	var videoBytes []byte
	videoBytes, err = ah.AvatarsUC.GetUserAvatar(userId)
	ah.answerMultipartAvatar(writer, videoBytes, http.StatusOK, err)
}
func (ah *AvatarsHandlers) answerEmpty(writer http.ResponseWriter, statusCode int, err error) {
	sent := ah.handlers.handleErrorNoContent(writer, err)
	if !sent {
		utils.CreateEmptyBodyAnswer(writer, statusCode)
	}
}

func (ah *AvatarsHandlers) answerMultipartAvatar(writer http.ResponseWriter, video []byte, statusCode int, err error) {
	sent := ah.handlers.handleErrorNoContent(writer, err)
	if !sent {
		utils.CreateAnswerMultipart(writer, statusCode, video)
	}
}
