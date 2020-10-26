package models

//easyjson:json
type WebsocketMessage struct {
	Type    int64   `json:"type"`
	Message Message `json:"message"`
}
