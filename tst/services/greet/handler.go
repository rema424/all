package greet

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HandleGreet ...s
func HandleGreet(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.String(http.StatusOK, fmt.Sprintf("Hello, %s!", name))
}
