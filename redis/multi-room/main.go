package main

import (
	"flag"
	"log"
	"net/http"
)

var redisPool = connectRedisPool()
var db = connectDB()
var e = createMux()

func init() {
	// Static
	e.Static("/js", "src/js")
	e.Static("/css", "src/css")

	// --------------------
	// Login
	// --------------------

	// Page
	e.GET("/login", loginPageHandler)
	e.GET("/signup", signupPageHandler)
	e.GET("/logout", logoutHandler)
	// API
	e.POST("/api/signup", signupExecHandler)
	e.POST("/api/login", loginExecHandler)

	// --------------------
	// Application
	// --------------------

	auth := e.Group("", isLoggedInHandler)

	// Page
	auth.GET("/", roomPageHandler)
	auth.GET("/rooms", roomPageHandler)
	auth.GET("/rooms/:roomID", roomPageHandler)
	// WebSocket
	auth.GET("/ws/rooms/:roomID", roomWebSocketHandler)
}

func main() {
	port := flag.String("port", "3000", "アプリケーションのアドレス")
	flag.Parse()

	http.Handle("/", e)
	log.Printf("server started at %s\n", *port)
	log.Fatal(http.ListenAndServe("localhost:"+*port, nil))
}
