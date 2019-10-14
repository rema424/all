package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"single-server/handler"
	"single-server/infra"
)

var e = infra.CreateMux()

func init() {

	e.GET("/", handler.TopHandler)
	e.GET("/signup", handler.SignUpPage)
	e.GET("/signin", handler.LoginPage)
	e.GET("/rooms", handler.RoomIndexPage)
	e.GET("/rooms/:roomID", handler.RoomShowPage2)
	e.GET("/rooms/:roomID/:userID", handler.RoomShowPage)
	e.GET("/rooms/:roomID/:userID/socket", handler.RoomShowWebSocket)
}

func main() {
	port := flag.String("port", "8080", "アプリケーションのポート")
	flag.Parse()

	http.Handle("/", e)
	log.Printf("Listening on port %s", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}
