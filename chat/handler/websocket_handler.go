package handler

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// ConnectWebSocket ...
func ConnectWebSocket(c echo.Context) error {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	wsconn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
	}
	defer wsconn.Close()

	for {
		// Write
		err := wsconn.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			log.Println(err)
		}

		// Read
		_, msg, err := wsconn.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(msg)
	}
}
