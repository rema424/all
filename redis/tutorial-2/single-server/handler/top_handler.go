package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HelloHandler ...
func HelloHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}

// TopHandler ...
func TopHandler(c echo.Context) error {
	return render(c, "top.html", map[string]interface{}{})
}
