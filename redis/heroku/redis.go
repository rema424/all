package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	conn := rr.pool.Get()
	defer conn.Close()
	psc := redis.PubSubConn{Conn: conn}
	psc.Subscribe(Channel)
	go rr.connHandler()
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			l.WithField("message", string(v.Data)).Info("Redis Message Reveive")
			if _, err := validateMessage(v.Data); err != nil {
				l.WithField("err", err).Error("Error unmarshalling message from Redis")
				continue
			}
			rr.broadcast(v.Data)
		case redis.Subscription:
			l.WithFields(logrus.Fields{
				"kind":  v.Kind,
				"count": v.Count,
			}).Println("Redis Subscription Reveived")
		case error:
			return errors.Wrap(v, "Error while subscribed to Redis channel")
		default:
			l.WithField("v", v).Info("Unknown Redis reveive during subscription")
		}
	}
}

func (rr *redisReceiver) broadcast(msg []byte) {
	rr.messages <- msg
}

func (rr *redisReceiver) register(conn *websocket.Conn) {
	rr.newConnections <- conn
}

func (rr *redisReceiver) deRegister(conn *websocket.Conn) {
	rr.rmConnections <- conn
}

func (rr *redisReceiver) connHandler() {
	conns := make([]*websocket.Conn, 0)
	for {
		select {
		case msg := <-rr.messages:
			for _, conn := range conns {
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.WithFields(logrus.Fields{
						"data": msg,
						"err":  err,
						"conn": conn,
					}).Error("Error writting data to connection! Closing and removing Connection")
					conns = removeConn(conns, conn)
				}
			}
		case conn := <-rr.newConnections:
			conns = append(conns, conn)
		case conn := <-rr.rmConnections:
			conns = removeConn(conns, conn)
		}
	}
}

func removeConn(conns []*websocket.Conn, remove *websocket.Conn) []*websocket.Conn {
	var i int
	var found bool
	for i = 0; i < len(conns); i++ {
		if conns[i] == remove {
			found = true
			break
		}
	}
	if !found {
		fmt.Printf("conns: %#v\nconn: %#v\n", conns, remove)
		panic("Conn not found")
	}
	copy(conns[i:], conns[i+1:]) // shift down
	conns[len(conns)-1] = nil    // nil last element
	return conns[:len(conns)-1]  // truncate slice
}
