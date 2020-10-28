package utils

import (
	"fmt"
	"github.com/google/logger"
	json "github.com/mailru/easyjson"
	masterhub_models "github.com/technoZoomers/MasterHubBackend/models"
	"net/http"
)

func writeData (w http.ResponseWriter, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		writeBytesError := fmt.Errorf("error writing video bytes to response: %v", err.Error())
		logger.Error(writeBytesError)
	}
}

func setHeaders(w http.ResponseWriter, statusCode int, contentType string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Access-Control-Allow-Origin", "Content-Type")
	if contentType != "" {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("content-type", contentType)
	}
}

func createAnswerJson(w http.ResponseWriter, statusCode int, data []byte) {
	setHeaders(w, statusCode, "application/json")
	writeData(w, data)
}

func CreateAnswerMultipart(w http.ResponseWriter, statusCode int, data []byte) {
	setHeaders(w, statusCode, "multipart/form-data;boundary=1")
	writeData(w, data)
}

func CreateEmptyBodyAnswer(writer http.ResponseWriter, statusCode int) {
	setHeaders(writer, statusCode, "")
}

func CreateErrorAnswerJson(writer http.ResponseWriter, statusCode int, error masterhub_models.RequestError) {
	marshalledError, err := json.Marshal(error)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Error(jsonError)
	}
	createAnswerJson(writer, statusCode, marshalledError)
}

func CreateAnswerLanguagesJson(writer http.ResponseWriter, statusCode int, languages masterhub_models.Languages) {
	marshalledLanguages, err := json.Marshal(languages)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Error(jsonError)
	}
	createAnswerJson(writer, statusCode, marshalledLanguages)
}

func CreateAnswerThemesJson(writer http.ResponseWriter, statusCode int, themes masterhub_models.Themes) {
	marshalledThemes, err := json.Marshal(themes)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Error(jsonError)
	}
	createAnswerJson(writer, statusCode, marshalledThemes)
}

func CreateAnswerThemeJson(writer http.ResponseWriter, statusCode int, theme masterhub_models.Theme) {
	marshalledTheme, err := json.Marshal(theme)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Error(jsonError)
	}
	createAnswerJson(writer, statusCode, marshalledTheme)
}

func CreateAnswerMasterJson(writer http.ResponseWriter, statusCode int, master masterhub_models.Master) {
	marshalledMaster, err := json.Marshal(master)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Error(jsonError)
	}
	createAnswerJson(writer, statusCode, marshalledMaster)
}

func CreateAnswerMasterFullJson(writer http.ResponseWriter, statusCode int, master masterhub_models.MasterFull) {
	marshalledMaster, err := json.Marshal(master)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Error(jsonError)
	}
	createAnswerJson(writer, statusCode, marshalledMaster)
}

func CreateAnswerMastersJson(writer http.ResponseWriter, statusCode int, masters masterhub_models.Masters) {
	marshalledMasters, err := json.Marshal(masters)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Error(jsonError)
	}
	createAnswerJson(writer, statusCode, marshalledMasters)
}

func CreateAnswerVideoDataJson(writer http.ResponseWriter, statusCode int, videoData masterhub_models.VideoData) {
	marshalledVideoData, err := json.Marshal(videoData)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Error(jsonError)
	}
	createAnswerJson(writer, statusCode, marshalledVideoData)
}

func CreateAnswerVideosDataJson(writer http.ResponseWriter, statusCode int, videosData masterhub_models.VideosData) {
	marshalledVideosData, err := json.Marshal(videosData)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Error(jsonError)
	}
	createAnswerJson(writer, statusCode, marshalledVideosData)
}

func CreateAnswerChatJson(writer http.ResponseWriter, statusCode int, chat masterhub_models.Chat) {
	marshalledChat, err := json.Marshal(chat)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Error(jsonError)
	}
	createAnswerJson(writer, statusCode, marshalledChat)
}

func CreateAnswerChatsJson(writer http.ResponseWriter, statusCode int, chats masterhub_models.Chats) {
	marshalledChats, err := json.Marshal(chats)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Error(jsonError)
	}
	createAnswerJson(writer, statusCode, marshalledChats)
}

func CreateAnswerMessagesJson(writer http.ResponseWriter, statusCode int, messages masterhub_models.Messages) {
	marshalledMessages, err := json.Marshal(messages)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Error(jsonError)
	}
	createAnswerJson(writer, statusCode, marshalledMessages)
}

func CreateAnswerUserJson(writer http.ResponseWriter, statusCode int, user masterhub_models.User) {
	marshalledUser, err := json.Marshal(user)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Error(jsonError)
	}
	createAnswerJson(writer, statusCode, marshalledUser)
}

func CreateAnswerStudentJson(writer http.ResponseWriter, statusCode int, student masterhub_models.Student) {
	marshalledStudent, err := json.Marshal(student)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Error(jsonError)
	}
	createAnswerJson(writer, statusCode, marshalledStudent)
}

func CreateAnswerStudentFullJson(writer http.ResponseWriter, statusCode int, student masterhub_models.StudentFull) {
	marshalledStudent, err := json.Marshal(student)
	if err != nil {
		jsonError := fmt.Errorf("error marshalling json: %v", err.Error())
		logger.Error(jsonError)
	}
	createAnswerJson(writer, statusCode, marshalledStudent)
}