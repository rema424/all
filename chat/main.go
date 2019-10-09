package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func main() {
	redisTutorial_1()
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
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}
