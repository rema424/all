package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"myproject/appengine/go-clean/controller"
	"myproject/domain/user"
	"myproject/gateway"
	"myproject/infra/mux"
	"myproject/infra/mysql"
	"myproject/infra/sqlxx"

	"github.com/labstack/echo/v4"
)

var (
	e      *echo.Echo = mux.CreateMux()
	accssr *sqlxx.Accessor
)

func init() {
	c := mysql.Config{
		Host:                 os.Getenv("DB_HOST"),
		Port:                 os.Getenv("DB_PORT"),
		User:                 os.Getenv("DB_USER"),
		DBName:               os.Getenv("DB_NAME"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		InterpolateParams:    true,
		AllowNatevePasswords: true,
		ParseTime:            true,
		MaxOpenConns:         -1, // use default value
		MaxIdleConns:         -1, // use default value
		ConnMaxLifetime:      -1, // use default value
	}

	dbx, err := mysql.Open(c)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("db opened successfully")

	accssr, err = sqlxx.Open(dbx)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("accessor opened successfully")

	ug := gateway.NewUserGateway(accssr)
	ui := user.NewInteractor(ug)
	uc := controller.NewUserController(ui)

	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "Hello, World.") })
	e.POST("/users", uc.Register)
}

func main() {
	http.Handle("/", e)
	defer accssr.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
