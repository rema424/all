package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/appengine"
)

var e = createMux()

func init() {
	// routing
	e.GET("/", helloHandler)
	e.GET("/worker", workerHandler)
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

func workerHandler(c echo.Context) error {
	tasks := 20
	msgs := make(chan string, tasks)

	for i := 1; i <= tasks; i++ {
		msgs <- fmt.Sprintf("task-%d", i)
	}

	var wg sync.WaitGroup
	workers := 3
	wg.Add(tasks)
	for i := 1; i <= workers; i++ {
		go worker(i, msgs, &wg)
	}
	close(msgs)
	wg.Wait()
	return nil
}

func worker(id int, msgs chan string, wg *sync.WaitGroup) {
	for msg := range msgs {
		task(id, msg)
		wg.Done()
	}
}

func task(workerID int, msg string) {
	s := rand.Intn(5) + 2
	time.Sleep(time.Second * time.Duration(s))
	fmt.Printf("workerID: %d, msg: %s, sleep: %d\n", workerID, msg, s)
}
