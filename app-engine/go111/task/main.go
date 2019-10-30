package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/appengine"
)

var e = createMux()

func init() {
	// routing
	e.GET("/", helloHandler)
	e.GET("/enqueue", enqueueHandler)
	e.GET("/dequeue", dequeueHandler)
}

func main() {
	http.Handle("/", e)

	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8080"
	// 	log.Printf("Defaulting to port %s", port)
	// }

	// log.Printf("Listening on port %s", port)
	// log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

	appengine.Main()
}

func createMux() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	return e
}

func helloHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}

func enqueueHandler(c echo.Context) error {
	return nil
}

func dequeueHandler(c echo.Context) error {
	return nil
}
