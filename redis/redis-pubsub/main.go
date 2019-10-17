package main

import (
	"os"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
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
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", os.Getenv("REDIS_URL")) },
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

func (c *Cache) newUser(socket *websocket.Conn, id string) *User {
	u := &User{
		ID:     id,
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

}
