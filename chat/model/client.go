package model

import (
	"github.com/gorilla/websocket"
)

// Client ...
type Client struct {
	hub     *Hub
	conn    *websocket.Conn
	user    *User
	message chan *Message
}
