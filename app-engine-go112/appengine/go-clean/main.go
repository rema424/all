// [START gae_go111_app]

// Sample helloworld is an App Engine app.
package main

// [START import]
import (
	"fmt"
	"log"
	"net/http"
	"os"

	"myproject/appengine/go-clean/controller"
	"myproject/appengine/go111/handler"
	"myproject/domain/user"
	"myproject/gateway"
	"myproject/infra/mux"
	"myproject/infra/mysql"
)

// [END import]
// [START main_func]

var e = mux.CreateMux()

func init() {
	ug := gateway.NewUserGateway()
	ui := user.NewInteractor(ug)
	uc := controller.NewUserController(ui)
	e.GET("/", handler.HandleHello)
	e.POST("/users", uc.Register)
}

func main() {
	defer mysql.Close()
	http.Handle("/", e)

	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	// [END setting_port]
}

// [END main_func]

// [START indexHandler]

// indexHandler responds to requests with our greeting.

// [END indexHandler]
// [END gae_go111_app]
