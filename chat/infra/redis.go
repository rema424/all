package infra

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

// ConnectRedis ...
func ConnectRedis() redis.Conn {
	var (
		host = mustGetenv("REDIS_HOST")
		port = mustGetenv("REDIS_PORT")
	)

	conn, err := redis.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("redis connected successfully")

	return conn
}
