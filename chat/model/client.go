package model

import (
	"time"

	"github.com/gorilla/websocket"
)

// WSClient ...
type WSClient struct {
	hub     *Hub
	wsconn  *websocket.Conn
	user    *User
	message chan *Message
}

const (
	pongWait       = 3 * time.Second
	pingPeriod     = 2 * time.Second
	maxMessageSize = 512
)

// // Upgrader ...
// var Upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }
