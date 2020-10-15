package models

type RequestError struct {
	Message string `json:"message"`
}

func CreateMessage(message string) RequestError {
	return RequestError{Message: message}
}
