package infra

import (
	"fmt"
	"log"
	"net/http"
)

var e = createMux()
var db = connectDB()
var rds = connectRedis()

// Run ...
func Run() {
	http.Handle("/", e)

	port := mustGetenv("APP_PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
