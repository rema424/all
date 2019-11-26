package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"tst/services/greet"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)

	handler := NewHandler(NewEcho())
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
	if err != nil {
		log.Fatal(err)
	}
}

// NewHandler ...
func NewHandler(e *echo.Echo) http.Handler {
	e.GET("/greet", greet.HandleGreet)
	return e
}

// NewEcho .
func NewEcho() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

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

// CreateCustomValidator .
func CreateCustomValidator(m map[string]validator.Func) *CustomValidator {
	v := validator.New()
	for key, val := range m {
		v.RegisterValidation(key, val)
	}
	return &CustomValidator{validator: v}
}

func sampleValidationRule(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	if len(val) > 30 {
		return true
	}
	return false
}

var reJapaneseZip = regexp.MustCompile(`^[\d]{3}-[\d]{4}$`)

// IsJapaneseZip は日本の郵便番号の形式をチェックします。
func IsJapaneseZip(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	// go1.11ではロック回避のためにCopy()して使う。go1.12からは必要ない。
	// ref: https://golang.org/doc/go1.12#regexp
	return reJapaneseZip.Copy().Match([]byte(val))
}
