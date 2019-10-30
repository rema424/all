package main

import (
	"encoding/json"
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

// Message ...
type Message struct {
	TaskName string `json:"taskName"`
	Content  string `json:"content"`
}

// Numbers ...
type Numbers []int

func enqueueHandler(c echo.Context) error {
	fmt.Println("enqueueHandler called")

	ctx := appengine.NewContext(c.Request())
	taskname := fmt.Sprintf("example-task-%d", time.Now().Unix())
	// msg := Message{taskname, "Hello from client"}
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8}
	// b, err := json.Marshal(msg)
	b, err := json.Marshal(nums)
	if err != nil {
		fmt.Println("json.Marshal() error:", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	t := taskqueue.Task{
		Path:    "/dequeue",
		Payload: b,
		Name:    taskname,
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
	// m := Message{}
	nums := []int{}
	err = json.Unmarshal(b, &nums)
	if err != nil {
		fmt.Println("json.Unmarshal() error:", err)
		return c.JSON(http.StatusInternalServerError, "わざと失敗")
	}

	for _, n := range nums {
		fmt.Println(n)
	}

	fmt.Println(nums)
	return c.JSON(http.StatusOK, nums)
}
