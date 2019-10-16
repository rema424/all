package main

import (
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("cmd", "go-websocket-chat-demo")
var waitTimeout = time.Minute * 10
var rr redisReceiver
var rw redisWriter

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.WithField("PORT", port).Fatal("$PORT must be set")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.WithField("REDIS_URL", redisURL).Fatal("$REDIS_URL must be set")
	}

	redisPool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", redisURL) },
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	// if err != nil {
	// 	log.WithField("url", redisURL).Fatal("Unable to create Redis pool")
	// }

	rr = newRedisReceiver(redisPool)
	rw = newRedisWriter(redisPool)

}
