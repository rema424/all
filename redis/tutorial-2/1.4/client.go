package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// clientはチャットを行なっている1人のユーザーを表します
type client struct {
	// socketはこのクライアントのためのWebSocketです
	socket *websocket.Conn
	// sendはメッセージが送られるチャネルです
	send chan []byte
	// roomはこのクライアントが参加しているチャットルームです
	room *room
}

func (c *client) read() {
	log.Println("client read started")
	defer log.Println("client read finished")

	for {
		log.Println("client read loop")
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	log.Println("client write started")
	defer log.Println("client write finished")

	for msg := range c.send {
		log.Println("client write loop")
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
