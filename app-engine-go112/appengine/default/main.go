// [START gae_go111_app]

// Sample helloworld is an App Engine app.
package main

// [START import]
import (
	"fmt"
	"log"
	"net/http"
	"os"

	"myproject/appengine/default/controller"
	"myproject/infra/mux"
)

// [END import]
// [START main_func]

var app = mux.CreateMux()

func init() {
	app.GET("/", controller.HandleHello)
}

func main() {
	http.Handle("/", app)

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
