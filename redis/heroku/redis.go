package main

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
)

// Channel ...
const Channel = "cnat"

var waitingMessage []byte
var availableMessage []byte
var waitSleep = time.Second * 10

func init() {
	var err error
	waitingMessage, err = json.Marshal(message{
		Handle: "system",
		Text:   "Waiting for redis to be available. Messaging won't work until redis is available",
	})
	if err != nil {
		panic(err)
	}
	availableMessage, err = json.Marshal(message{
		Handle: "system",
		Text:   "Redis is now available & messaging is now possible",
	})
	if err != nil {
		panic(err)
	}
}

type redisReceiver struct {
	pool           *redis.Pool
	messages       chan []byte
	newConnections chan *websocket.Conn
	rmConnections  chan *websocket.Conn
}

func newRedisReceiver(pool *redis.Pool) redisReceiver {
	return redisReceiver{
		pool:           pool,
		messages:       make(chan []byte, 1000),
		newConnections: make(chan *websocket.Conn),
		rmConnections:  make(chan *websocket.Conn),
	}
}

func (rr *redisReceiver) wait(_ time.Time) error {
	rr.broadcast(waitingMessage)
	time.Sleep(waitSleep)
	return nil
}

func (rr *redisReceiver) run() error {
	l := log.WithField("channel", Channel)
	conn := rr.pool
}
