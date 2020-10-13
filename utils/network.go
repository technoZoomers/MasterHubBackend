package utils

import (
	"github.com/google/logger"
	json "github.com/mailru/easyjson"
	masterhub_models "github.com/technoZoomers/MasterHubBackend/models"
	"net/http"
)

func createAnswerJson(w http.ResponseWriter, statusCode int, data []byte) {
	w.WriteHeader(statusCode)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	_, err := w.Write(data)
	if err != nil {
		logger.Errorf("Error writing answer: %v", err)
	}
}

func createEmptyBodyAnswerJson(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
}

func CreateEmptyBodyAnswerJson(writer http.ResponseWriter, statusCode int) {
	createEmptyBodyAnswerJson(writer, statusCode)
}

func CreateErrorAnswerJson(writer http.ResponseWriter, statusCode int, error masterhub_models.RequestError) {
	marshalledError, err := json.Marshal(error)
	if err != nil {
		logger.Errorf("Error marhalling json: %v", err)
	}
	createAnswerJson(writer, statusCode, marshalledError)
}

func CreateAnswerLanguagesJson(writer http.ResponseWriter, statusCode int, languages masterhub_models.Languages) {
	marshalledLanguages, err := json.Marshal(languages)
	if err != nil {
		logger.Errorf("Error marhalling json: %v", err)
	}
	createAnswerJson(writer, statusCode, marshalledLanguages)
}

func CreateAnswerThemesJson(writer http.ResponseWriter, statusCode int, themes masterhub_models.Themes) {
	marshalledThemes, err := json.Marshal(themes)
	if err != nil {
		logger.Errorf("Error marhalling json: %v", err)
	}
	createAnswerJson(writer, statusCode, marshalledThemes)
}

func CreateAnswerThemeJson(writer http.ResponseWriter, statusCode int, theme masterhub_models.Theme) {
	marshalledTheme, err := json.Marshal(theme)
	if err != nil {
		logger.Errorf("Error marhalling json: %v", err)
	}
	createAnswerJson(writer, statusCode, marshalledTheme)
}

func CreateAnswerMasterJson(writer http.ResponseWriter, statusCode int, master masterhub_models.Master) {
	marshalledMaster, err := json.Marshal(master)
	if err != nil {
		logger.Errorf("Error marhalling json: %v", err)
	}
	createAnswerJson(writer, statusCode, marshalledMaster)
}

func CreateAnswerVideoDataJson(writer http.ResponseWriter, statusCode int, videoData masterhub_models.VideoData) {
	marshalledVideoData, err := json.Marshal(videoData)
	if err != nil {
		logger.Errorf("Error marhalling json: %v", err)
	}
	createAnswerJson(writer, statusCode, marshalledVideoData)
}

func CreateAnswerVideosDataJson(writer http.ResponseWriter, statusCode int, videosData masterhub_models.VideosData) {
	marshalledVideosData, err := json.Marshal(videosData)
	if err != nil {
		logger.Errorf("Error marhalling json: %v", err)
	}
	createAnswerJson(writer, statusCode, marshalledVideosData)
}