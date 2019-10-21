package main

import (
	"context"

	"github.com/labstack/echo/v4"
)

// GetDBx ...
func GetDBx(c echo.Context) *DBx {
	if obj := c.Get("dbx"); obj != nil {
		if dbx, ok := obj.(*DBx); ok {
			return dbx
		}
	}
	dbx := NewDBx(context.Background(), db)
	SetDBx(c, dbx)
	return dbx
}

// SetDBx ...
func SetDBx(c echo.Context, dbx *DBx) {
	c.Set("dbx", dbx)
}

// ReadSessID ...
func ReadSessID(c echo.Context) string {
	cookie, err := c.Cookie("sessid")
	if err != nil {
		return ""
	}
	return cookie.Value
}
