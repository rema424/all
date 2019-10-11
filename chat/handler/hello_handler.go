package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HelloHandler ...
func HelloHandler(c echo.Context) error {
	// for i := 0; i < 10; i++ {
	// 	// lib.UUID()
	// }
	return c.String(http.StatusOK, "Hello World!")
}
