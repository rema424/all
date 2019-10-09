package infra

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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
		// lib.UUID()
	}
	return c.String(http.StatusOK, "Hello World!")
}
