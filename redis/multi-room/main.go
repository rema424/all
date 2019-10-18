package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
)

var redisPool = &redis.Pool{
	MaxIdle:     3,
	IdleTimeout: 240 * time.Second,
	Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", ":6379") },
	TestOnBorrow: func(c redis.Conn, t time.Time) error {
		if time.Since(t) < time.Minute {
			return nil
		}
		_, err := c.Do("PING")
		return err
	},
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	port := flag.String("port", ":3000", "アプリケーションのアドレス")
	flag.Parse()

	http.HandleFunc("/rooms/", roomHandler)

	log.Printf("server started at %s\n", *port)
	log.Fatal(http.ListenAndServe("localhost:"+*port, nil))
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

// /rooms/{roomID}
func roomHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	// ルームを準備する
	roomID, _ := strconv.Atoi(segs[2])
	var room *Room
	if val, ok := roomCache[roomID]; ok {
		fmt.Println("キャッシュから部屋インスタンスを取得しました。")
		room = val
	} else {
		fmt.Printf("部屋ID: %d のインスタンスを新たに作成しました。\n", roomID)
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
				fmt.Printf("Redisによる部屋ID: %d の購読を開始します。\n", room.roomID)
				defer fmt.Printf("Redisによる部屋ID: %d の購読を終了します。\n", room.roomID)

				psc.Subscribe(room.roomID)
				for {
					switch v := psc.Receive().(type) {
					case redis.Message:
						fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
					case redis.Subscription:
						fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
					case error:
						fmt.Println("close psc")

						close(room.doneCh)
						return
					}
				}
			}()

			fmt.Printf("部屋ID: %d の監視を開始します。\n", roomID)
			defer fmt.Printf("部屋ID: %d の監視を終了。\n", roomID)
			for {
				select {
				case <-room.doneCh:
					fmt.Printf("部屋ID: %d に閉鎖の通知が届きました。\n", roomID)
					delete(roomCache, roomID)
					for client := range room.Clients {
						client.socket.Close()
					}
					return
				case client := <-room.newClientCh:
					fmt.Printf("部屋ID: %d に入室の通知が届きました。\n", roomID)
					room.Clients[client] = true
				case client := <-room.rmClientCh:
					fmt.Printf("部屋ID: %d に退室の通知が届きました。\n", roomID)
					delete(room.Clients, client)
				case msg := <-room.msgCh:
					fmt.Printf("部屋ID: %d へメッセージを送信します。\n", roomID)
					if c := redisPool.Get(); c != nil {
						c.Do("PUBLISH", room.roomID, msg)
						c.Close()
					}
				}
			}
		}()
	})

	// クライアントを準備する
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrader error %s\n", err.Error())
		return
	}
	client := &Client{socket: socket}

	// registerClient()
	func() {
		fmt.Printf("部屋ID: %d へ入室します。\n", roomID)
		room.newClientCh <- client
	}()

	// deRegisterClient()
	defer func() {
		fmt.Printf("部屋ID: %d から退室します。\n", roomID)
		room.rmClientCh <- client
	}()

	// messageServerToClients()
	go func() {
		fmt.Printf("部屋ID: %d: ブラウザへメッセージを送信する準備をします。\n", roomID)
		for {

		}
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
}
