package models

import "github.com/gorilla/websocket"

type WebsocketConnection struct {
	UserId int64
	Connection *websocket.Conn
}
