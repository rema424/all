package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HandleHello ...
func HandleHello(c echo.Context) error {
	// ctx := c.Request().Context()
	// ctx
	return c.String(http.StatusOK, "Hello, World! go111")
}
