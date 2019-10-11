package main

import (
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("cmd", "go-websocket-chat-demo")

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	redisPool := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   0,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", "localhost:6379") },
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

type redisReceiver struct {
	pool *redis.Pool

	messages       chan []byte
	newConnections chan *websocket.Conn
	rmConnections  chan *websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type message struct {
	Handle string `json:"handle"`
	Text   string `json:"text"`
}

// ChannelName to use with redis
const ChannelName = "chat"

var waitingMessage []byte
var availableMessage []byte
var waitSleep = 10 * time.Second

func init() {
	// var err error
	// waitingMessage, err = json.Marshal()
}
