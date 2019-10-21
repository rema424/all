package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Room ...
type Room struct {
	roomID      int
	psc         *redis.PubSubConn
	once        sync.Once
	doneCh      chan struct{}
	newClientCh chan *Client
	rmClientCh  chan *Client
	msgCh       chan []byte
	Clients     map[*Client]bool
}

// Client ...
type Client struct {
	socket *websocket.Conn
}

var roomCache = map[int]*Room{}

// --------------------------------------------------
// roomIndexPageHandler
// --------------------------------------------------

func roomIndexPageHandler(c echo.Context) error {
	if c.Request().URL.Path == "/" {
		c.Redirect(http.StatusPermanentRedirect, "/rooms")
	}
	return render(c, "chat.html", map[string]interface{}{})
}

// --------------------------------------------------
// roomShowPageHandler
// --------------------------------------------------

func roomShowPageHandler(c echo.Context) error {
	roomID := c.Param("roomID")
	return render(c, "chat.html", map[string]interface{}{
		"RoomID": roomID,
	})
}

// --------------------------------------------------
// roomWebSocektHandler
// --------------------------------------------------

func roomWebSocketHandler(c echo.Context) error {
	defer fmt.Println("リクエストの処理を終了します。")
	// ルームを準備する
	roomID, _ := strconv.Atoi(c.Param("roomID"))
	var room *Room
	if val, ok := roomCache[roomID]; ok {
		fmt.Printf("部屋ID: %d: キャッシュから部屋インスタンスを取得しました。\n", roomID)
		room = val
	} else {
		fmt.Printf("部屋ID: %d: 部屋インスタンスを新たに作成しました。\n", roomID)
		room = &Room{
			roomID:      roomID,
			doneCh:      make(chan struct{}),
			newClientCh: make(chan *Client),
			rmClientCh:  make(chan *Client),
			msgCh:       make(chan []byte),
			Clients:     make(map[*Client]bool),
		}
		roomCache[roomID] = room
	}

	// run()
	room.once.Do(func() {
		go func() {
			c := redisPool.Get()
			defer c.Close()
			psc := redis.PubSubConn{Conn: c}
			defer psc.Close()

			go func() {
				fmt.Printf("部屋ID: %d: Redisの購読を開始します。\n", room.roomID)
				defer fmt.Printf("部屋ID: %d: Redisの購読を終了します。\n", room.roomID)

				psc.Subscribe(room.roomID)
				for {
					switch v := psc.Receive().(type) {
					case redis.Message:
						fmt.Printf("部屋ID: %d: Redisからメッセージが届きました。クライアントに転送します。\n", room.roomID)
						fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
						for client := range room.Clients {
							if err := client.socket.WriteMessage(websocket.TextMessage, v.Data); err != nil {
								fmt.Printf("WebSocketでメッセージ送信に失敗しました。 %v\n", err)
							} else {
								fmt.Printf("WebSocketでメッセージを送信しました。")
							}
						}
					case redis.Subscription:
						fmt.Printf("部屋ID: %d: Redisから購読開始通知が届きました。\n", room.roomID)
						fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
					case error:
						fmt.Printf("部屋ID: %d: Redisの購読エラーです。%v\n", room.roomID, v.Error())
						fmt.Println("close psc")
						close(room.doneCh)
						return
					}

				}
			}()

			fmt.Printf("部屋ID: %d: 部屋の監視を開始します。\n", roomID)
			defer fmt.Printf("部屋ID: %d: 部屋の監視を終了します。\n", roomID)
			for {
				select {
				case <-room.doneCh:
					fmt.Printf("部屋ID: %d: 部屋の閉鎖通知が届きました。\n", roomID)
					delete(roomCache, roomID)
					for client := range room.Clients {
						client.socket.Close()
					}
					return
				case client := <-room.newClientCh:
					fmt.Printf("部屋ID: %d: 入室の通知が届きました。\n", roomID)
					room.Clients[client] = true
				case client := <-room.rmClientCh:
					fmt.Printf("部屋ID: %d: 退室の通知が届きました。\n", roomID)
					delete(room.Clients, client)
				case msg := <-room.msgCh:
					fmt.Printf("部屋ID: %d: 部屋にメッセージが届きました。Redisに転送します。\n", roomID)
					if c := redisPool.Get(); c != nil {
						c.Do("PUBLISH", room.roomID, msg)
						c.Close()
					}
				}
			}
		}()
	})

	// クライアントを準備する
	socket, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("upgrader error %s\n", err.Error())
		return nil
	}
	defer socket.Close()

	client := &Client{socket: socket}

	// registerClient()
	func() {
		fmt.Printf("部屋ID: %d: 入室します。\n", roomID)
		room.newClientCh <- client
	}()

	// deRegisterClient()
	defer func() {
		fmt.Printf("部屋ID: %d: 退室します。\n", roomID)
		room.rmClientCh <- client
	}()

	// messageServerToClients()
	go func() {
		fmt.Printf("部屋ID: %d: ブラウザへメッセージを送信する準備をします。\n", roomID)
		// for {

		// }
	}()

	// messageClientToServer()
	func() {
		fmt.Printf("部屋ID: %d: ブラウザからメッセージを受信する準備をします。\n", roomID)
		for {
			if _, msg, err := client.socket.ReadMessage(); err == nil {
				fmt.Printf("部屋ID: %d: ブラウザからメッセージを受信しました。\n", roomID)
				room.msgCh <- msg
			} else {
				break
			}
		}
		client.socket.Close()
		fmt.Println("ブラウザとの接続が切れました。")
	}()
	return nil
}
