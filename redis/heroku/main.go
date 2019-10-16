package main

import (
	"net/url"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/gomodule/redigo/redis"
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

	u, err := url.Parse(redisURL)
	if err != nil {
		panic(err)
	}

	var pwds []string
	pass, ok := u.User.Password()
	if ok {
		pwds = append(pwds, pass)
		u.User = url.UserPassword("", "")
	}

	redisPool := &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial:
	}
}
