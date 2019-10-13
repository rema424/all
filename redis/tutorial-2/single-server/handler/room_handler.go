package handler

import (
	"fmt"
	"strconv"

	"single-server/service"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// RoomIndexPage ...
func RoomIndexPage(c echo.Context) error {
	return nil
}

// RoomShowPage ...
func RoomShowPage(c echo.Context) error {
	// roomID, err := strconv.Atoi(c.Param("roomID"))
	// if err != nil {
	// 	fmt.Println("RoomIDの取得に失敗しました：", roomID, "-", err)
	// }
	// fmt.Println("RoomIDの取得に成功しました：", roomID)

	// userID, err := strconv.Atoi(c.Param("userID"))
	// if err != nil {
	// 	fmt.Println("userIDの取得に失敗しました：", userID, "-", err)
	// }
	// fmt.Println("userIDの取得に成功しました：", userID)

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
		fmt.Println("roomIDの取得に失敗しました：", roomID, "-", err)
	}
	fmt.Println("roomIDの取得に成功しました：", roomID)

	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		fmt.Println("userIDの取得に失敗しました：", userID, "-", err)
	}
	fmt.Println("userIDの取得に成功しました：", userID)

	socket, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println("WebSocket接続の取得に失敗しました：", socket, "-", err)
		return err
	}
	defer socket.Close()
	defer fmt.Println("ソケットを閉じました。")
	fmt.Println("WebSocket接続の取得に成功しました：")

	service.ConnectChatRoom(nil, socket, roomID, userID)
	return nil
}
