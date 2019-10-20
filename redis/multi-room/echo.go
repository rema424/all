package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
)

func createMux() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Validator = createCustomValidator()
	return e
}

// CustomValidator ...
type CustomValidator struct {
	validator *validator.Validate
}

// Validate ...
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func createCustomValidator() *CustomValidator {
	v := validator.New()
	v.RegisterValidation("myrule", myValidationRule)
	return &CustomValidator{validator: v}
}

func myValidationRule(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	if len(val) > 30 {
		return true
	}
	return false
}

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