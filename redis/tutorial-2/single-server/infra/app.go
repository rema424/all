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
	e.Static("/css", "asset/css")
	e.Static("/js", "asset/js")

	return e
}
