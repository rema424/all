package main

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

func main() {
	// redisTutorial_1()
	redisTutorial_2()
}

func redisTutorial_1() {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 1
	reply, err := conn.Do("get", "foo")
	if err != nil {
		fmt.Println("1-err", err)
	} else {
		fmt.Println("1-rep", reply)
	}

	// 2
	reply, err = redis.String(conn.Do("get", "foo"))
	if err != nil {
		fmt.Println("2-err", err)
	} else {
		fmt.Println("2-rep", reply)
	}

	// 3
	reply, err = conn.Do("set", "foo", "bar")
	if err != nil {
		fmt.Println("3-err", err)
	} else {
		fmt.Println("3-rep", reply)
	}

	// 4
	reply, err = conn.Do("get", "foo")
	if err != nil {
		fmt.Println("4-err", err)
	} else {
		fmt.Println("4-rep", reply)
	}

	// 5
	reply, err = redis.String(conn.Do("get", "foo"))
	if err != nil {
		fmt.Println("5-err", err)
	} else {
		fmt.Println("5-rep", reply)
	}

	// 6
	reply, err = conn.Do("del", "foo")
	if err != nil {
		fmt.Println("6-err", err)
	} else {
		fmt.Println("6-rep", reply)
	}

	// 7
	reply, err = conn.Do("get", "foo")
	if err != nil {
		fmt.Println("7-err", err)
	} else {
		fmt.Println("7-rep", reply)
	}
}

func redisTutorial_2() {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
		IdleTimeout: 4 * 60 * time.Second,
		MaxActive:   6,
		MaxIdle:     3,
	}

	printStatus := func(pool *redis.Pool) {
		s := pool.Stats()
		fmt.Printf("Active: %d, Idle: %d, InUse: %d\n", s.ActiveCount, s.IdleCount, s.ActiveCount-s.IdleCount)
	}

	printStatus(pool)

	fmt.Println("------")

	conns := make([]redis.Conn, 10)
	for i := 0; i < 10; i++ {
		conns[i] = pool.Get()
		fmt.Printf("%d: ", i)
		printStatus(pool)
	}

	fmt.Println("------")

	for i, conn := range conns {
		conn.Close()
		fmt.Printf("%d: ", i)
		printStatus(pool)
	}
}
