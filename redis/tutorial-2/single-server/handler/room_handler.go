package handler

import (
	"fmt"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// RoomIndexPage ...
func RoomIndexPage(c echo.Context) error {
	return nil
}

// RoomShowPage ...
func RoomShowPage(c echo.Context) error {
	roomID, err := strconv.Atoi(c.Param("roomID"))
	if err != nil {
		fmt.Println("RoomIDの取得に失敗しました：", roomID, "-", err)
	}
	fmt.Println("RoomIDの取得に成功しました：", roomID)
	return render(c, "room/show.html", map[string]interface{}{"Host": ":8080"})
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

// RoomShowWebSocket ...
func RoomShowWebSocket(c echo.Context) error {
	roomID, err := strconv.Atoi(c.Param("roomID"))
	if err != nil {
		fmt.Println("RoomIDの取得に失敗しました：", roomID, "-", err)
	}
	fmt.Println("RoomIDの取得に成功しました：", roomID)

	socket, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println("WebSocket接続の取得に失敗しました：", socket, "-", err)
		return err
	}
	defer socket.Close()
	fmt.Println("WebSocket接続の取得に成功しました：", socket)

	for {

	}
}
