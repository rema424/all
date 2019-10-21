package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func roomIndexPageHandler(c echo.Context) error {
	if c.Request().URL.Path == "/" {
		c.Redirect(http.StatusPermanentRedirect, "/rooms")
	}
	return render(c, "chat.html", map[string]interface{}{})
}
