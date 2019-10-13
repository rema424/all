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

	e.GET("/", handler.HelloHandler)
	e.GET("/rooms/:roomID", handler.RoomShowPage)
}

func main() {
	port := flag.String("port", "8080", "アプリケーションのポート")
	flag.Parse()

	http.Handle("/", e)
	log.Printf("Listening on port %s", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}
