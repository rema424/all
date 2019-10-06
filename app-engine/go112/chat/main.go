package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"chat/lib"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var e = createMux()

func main() {
	http.Handle("/", e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func createMux() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	// routing
	e.GET("/", helloHandler)

	return e
}

func helloHandler(c echo.Context) error {
	for i := 0; i < 10; i++ {
		lib.UUID()
	}
	return c.String(http.StatusOK, "Hello World!")
}
