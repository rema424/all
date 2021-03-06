package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

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

// RoomRecord ...
type RoomRecord struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// Message ...
type Message struct {
	ID        int       `db:"id"`
	RoomID    int       `db:"room_id"`
	UserID    int       `db:"user_id"`
	Body      string    `db:"body"`
	CreatedAt time.Time `db:"created_at"`
	User      *User     `db:"user"`
}

var roomCache = map[int]*Room{}

// --------------------------------------------------
// roomPageHandler
// --------------------------------------------------

func roomPageHandler(c echo.Context) error {
	if c.Request().URL.Path == "/" {
		return c.Redirect(http.StatusPermanentRedirect, "/rooms")
	}

	cookie, err := c.Cookie("sessid")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login")
	}
	sessid := cookie.Value

	dbx := GetDBx(c)

	var u User
	q := `
  select u.id as id, u.name as name
  from user as u
  inner join session as s on s.user_id = u.id
  where s.session_id = ?;`
	if err := dbx.Get(&u, q, sessid); err != nil {
		fmt.Println(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	var rooms []*RoomRecord
	q = "select * from room;"
	if err := dbx.Select(&rooms, q); err != nil {
		fmt.Println(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	roomID, _ := strconv.Atoi(c.Param("roomID"))
	if roomID == 0 {
		return render(c, "chat.html", map[string]interface{}{
			"User":  u,
			"Rooms": rooms,
		})
	}

	var msgs []*Message
	q = `
  select
    m.id as id,
    m.body as body,
    m.created_at as created_at,
    u.id as 'user.id',
    u.name as 'user.name'
  from message as m
  inner join user as u on u.id = m.user_id
  where m.room_id = ?;`
	if err := dbx.Select(&msgs, q, roomID); err != nil {
		fmt.Println(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	// return c.JSONPretty(200, msgs, "  ")

	return render(c, "chat.html", map[string]interface{}{
		"User":     u,
		"Rooms":    rooms,
		"RoomID":   roomID,
		"Messages": msgs,
	})
}

// --------------------------------------------------
// roomWebSocektHandler
// --------------------------------------------------

func roomWebSocketHandler(c echo.Context) error {
	defer fmt.Println("リクエストの処理を終了します。")

	// ユーザーを取得する
	cookie, _ := c.Cookie("sessid")
	if cookie == nil || cookie.Value == "" {
		return nil
	}
	sessid := cookie.Value
	dbx := GetDBx(c)
	var u User
	q := `
  select
    u.id as id,
    u.name as name
  from user as u
  inner join session as s on s.user_id = u.id
  where s.session_id = ?;`
	if err := dbx.Get(&u, q, sessid); err != nil {
		return nil
	}

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
				m := Message{
					Body:      string(msg),
					CreatedAt: time.Now(),
					User:      &u,
				}
				b, err := json.Marshal(m)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					room.msgCh <- b
				}
			} else {
				break
			}
		}
		client.socket.Close()
		fmt.Println("ブラウザとの接続が切れました。")
	}()
	return nil
}
