package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/appengine"
	"google.golang.org/appengine/taskqueue"
)

var e = createMux()

func init() {
	// routing
	e.GET("/", helloHandler)
	e.GET("/enqueue", enqueueHandler)
	e.POST("/dequeue", dequeueHandler)
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
	fmt.Println("enqueueHandler called")

	ctx := appengine.NewContext(c.Request())
	t := taskqueue.Task{
		Path:    "/dequeue",
		Payload: []byte("example"),
		Name:    fmt.Sprintf("example-task-%d", time.Now().Unix()),
		Delay:   1 * time.Second,
		RetryOptions: &taskqueue.RetryOptions{
			RetryLimit: 4,
			MinBackoff: 500 * time.Millisecond,
		},
	}

	if _, err := taskqueue.Add(ctx, &t, "queue-blue"); err != nil {
		fmt.Println("taskqueue.Add() error:", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	fmt.Println("task added successfully")
	return c.JSON(http.StatusOK, "task added successfully")
}

func dequeueHandler(c echo.Context) error {
	fmt.Println("dequeueHandler called")
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println("Body読み込みエラー:", err)
		return c.JSON(http.StatusInternalServerError, "わざと失敗")
	}

	fmt.Println(string(b))
	return c.JSON(http.StatusOK, string(b))
}
