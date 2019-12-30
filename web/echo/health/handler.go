package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleHealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"alive": true})
}
