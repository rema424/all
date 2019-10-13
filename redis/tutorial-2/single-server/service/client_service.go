package service

import (
	"fmt"
	"log"
	"single-server/infra"
	"single-server/model"
	"time"

	"github.com/gorilla/websocket"
)

const messageBufferSize = 256

type clientSupervisor struct {
	socket *websocket.Conn
	send   chan *model.Message
	room   *model.Room
	user   *model.User
}

// ConnectClient ...
func ConnectClient(db *infra.DB, socket *websocket.Conn, roomID, userID int) {
	room, err := GetRoomByID(db, roomID)
	if err != nil {
		fmt.Println("Roomの取得に失敗しました：", room, "-", err)
	}

	user, err := GetUserByID(db, userID)
	if err != nil {
		fmt.Println("Userの取得に失敗しました：", user, "-", err)
	}

	client := &clientSupervisor{
		socket: socket,
		send:   make(chan *model.Message, model.MessageBufferSize),
		room:   room,
		user:   user,
	}

	room.JoinCh <- client
	defer func() { room.LeaveCh <- client }()
	go client.Write()
	client.Read()

}

func (c *clientSupervisor) Read() {
	defer c.socket.Close()

	log.Println("client read started")
	defer log.Println("client read finished")

	for {
		log.Println("client read loop")
		var msg *model.Message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.User = c.user
			c.room.MessageCh <- msg
		} else {
			break
		}
	}
}

func (c *clientSupervisor) Write() {
	defer c.socket.Close()

	log.Println("client write started")
	defer log.Println("client write finished")

	for msg := range c.send {
		log.Println("client write loop")
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
}

// JoinRoom ...
// func JoinRoom(db *infra.DB, roomID int, socket *websocket.Conn, user *model.User) {
// 	room, err := GetRoomByID(db, roomID)
// 	if err != nil {
// 		fmt.Println("Roomの取得に失敗しました：", room, "-", err)
// 	}

// 	client := &model.Client{
// 		Socket: socket,
// 		Send:   make(chan *model.Message, model.MessageBufferSize),
// 		Room:   room,
// 		User:   user,
// 	}

// 	room.JoinCh <- client
// 	defer func() { room.LeaveCh <- client }()
// 	go client.Write()
// 	client.Read()
// }

// GetRoomByID ...
// func GetRoomByID(db *infra.DB, roomID int) (*model.Room, error) {
// 	return nil, nil
// }
