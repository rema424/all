package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// clientはチャットを行なっている1人のユーザーを表します
type client struct {
	// socketはこのクライアントのためのWebSocketです
	socket *websocket.Conn
	// sendはメッセージが送られるチャネルです
	send chan *message
	// roomはこのクライアントが参加しているチャットルームです
	room *room
	// userDataはユーザーに関する情報を保持します
	userData map[string]interface{}
}

func (c *client) read() {
	log.Println("client read started")
	defer log.Println("client read finished")

	for {
		log.Println("client read loop")
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
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
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
