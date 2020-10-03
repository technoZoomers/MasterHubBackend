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

func CreateAnswerBalanceJson(writer http.ResponseWriter, statusCode int, chats masterhub_models.Balance) {
	marshalledBalance, err := json.Marshal(chats)
	if err != nil {
		logger.Errorf("Error marhalling json: %v", err)
	}
	createAnswerJson(writer, statusCode, marshalledBalance)
}

func CreateAnswerTransactionsJson(writer http.ResponseWriter, statusCode int, chats masterhub_models.Transactions) {
	marshalledTransactions, err := json.Marshal(chats)
	if err != nil {
		logger.Errorf("Error marhalling json: %v", err)
	}
	createAnswerJson(writer, statusCode, marshalledTransactions)
}
