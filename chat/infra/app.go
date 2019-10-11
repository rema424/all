package infra

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CreateMux ...
func CreateMux() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	// routing
	// e.GET("/", handler.HelloHandler)

	return e
}
