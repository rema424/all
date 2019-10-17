package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// User ...
type User struct {
	ID     string
	socket *websocket.Conn
}

// Cache ...
type Cache struct {
	Users []*User
	sync.Mutex
}

// Message ...
type Message struct {
	DeliveryID string `json:"id"`
	Content    string `json:"content"`
}

var (
	cache     *Cache
	pubSub    *redis.PubSubConn
	redisPool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", ":6379") },
		// Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", os.Getenv("REDIS_URL")) },
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
)

func init() {
	cache = &Cache{
		Users: make([]*User, 0, 1),
	}
}

var serverAddress = ":3000"

// func (c *Cache) newUser(socket *websocket.Conn, id string) *User {
func (c *Cache) newUser(socket *websocket.Conn) *User {
	u := &User{
		// ID:     id,
		ID:     uuid.New().String(),
		socket: socket,
	}

	if err := pubSub.Subscribe(u.ID); err != nil {
		panic(err)
	}
	c.Lock()
	defer c.Unlock()

	c.Users = append(c.Users, u)
	return u
}

func main() {
	redisConn := redisPool.Get()
	defer redisConn.Close()

	pubSub = &redis.PubSubConn{Conn: redisConn}
	defer pubSub.Close()

	go deliverMessage()

	http.HandleFunc("/ws", wsHandler)

	log.Printf("server started at %s\n", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

func deliverMessage() {
	for {
		switch v := pubSub.Receive().(type) {
		case redis.Message:
			cache.findAndDeliver(v.Channel, string(v.Data))
		case redis.Subscription:
			log.Printf("subscription message: %s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			log.Printf("error pub/sub on connection, delivery has stopped")
			return
		}
	}
}

func (c *Cache) findAndDeliver(userID string, content string) {
	m := Message{
		Content: content,
	}

	for _, u := range c.Users {
		if u.ID == userID {
			if err := u.socket.WriteJSON(m); err != nil {
				log.Printf("error on message delivery through ws. e: %s\n", err)
			} else {
				log.Printf("user %s found at our store, message sent\n", userID)
			}
			return
		}
	}

	log.Printf("user %s not found at our store\n", userID)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrader error %s\n", err.Error())
		return
	}
	u := cache.newUser(socket)
	// u := cache.newUser(socket, r.FormValue("id"))
	log.Printf("user %s joined\n", u.ID)

	for {
		var m Message

		if err := u.socket.ReadJSON(&m); err != nil {
			log.Printf("error on ws. message %s\n", err.Error())
		}

		if c := redisPool.Get(); c != nil {
			c.Do("PUBLISH", m.DeliveryID, string(m.Content))
		}
	}
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
		room = val
	} else {
		room = &Room{roomID: roomID}
	}

	// run()
	room.once.Do(func() {
		go func() {
		loop:
			for {
				select {
				case <-room.doneCh:
					delete(roomCache, roomID)
					for client := range room.Clients {
						client.socket.Close()
					}
					break loop
				case client := <-room.newClientCh:
					room.Clients[client] = true
				case client := <-room.rmClientCh:
					delete(room.Clients, client)
				case msg := <-room.msgCh:
					if c := redisPool.Get(); c != nil {
						c.Do("PUBLISH", room.roomID, msg)
						c.Close()
					}
					fmt.Println(msg)
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
		room.newClientCh <- client
	}()

	// deRegisterClient()
	defer func() {
		room.rmClientCh <- client
	}()

	// messageServerToClients()
	go func() {
		for {

		}
	}()

	// messageClientToServer()
	func() {
		for {
			if _, msg, err := client.socket.ReadMessage(); err == nil {
				room.msgCh <- msg
			} else {
				break
			}
		}
		client.socket.Close()
	}()

	u := cache.newUser(socket)
	// u := cache.newUser(socket, r.FormValue("id"))
	log.Printf("user %s joined\n", u.ID)

	// ルームはクライアントを登録する

	for {
		var m Message

		if err := u.socket.ReadJSON(&m); err != nil {
			log.Printf("error on ws. message %s\n", err.Error())
		}

		if c := redisPool.Get(); c != nil {
			c.Do("PUBLISH", m.DeliveryID, string(m.Content))
		}
	}
}
